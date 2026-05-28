package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"alhikmah-attendance-api/internal/domain"
)

type reportService struct {
	reportRepo domain.ReportRepository
	cache      domain.ReportCacheRepository
}

func NewReportService(reportRepo domain.ReportRepository, studentRepo domain.StudentRepository, cacheRepo domain.ReportCacheRepository) domain.ReportService {
	// studentRepo is not used here but kept for interface compatibility
	return &reportService{
		reportRepo: reportRepo,
		cache:      cacheRepo,
	}
}

func (s *reportService) GetDailyReport(classID, dateStr string, forceRefresh bool, generatedBy string) (*domain.DailyReport, error) {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	if forceRefresh {
		s.cache.Delete("harian", classID, dateStr, dateStr)
	} else {
		cached, err := s.cache.Get("harian", classID, dateStr, dateStr)
		if err == nil && cached != nil {
			var report domain.DailyReport
			if err := json.Unmarshal(cached, &report); err == nil {
				return &report, nil
			}
		}
	}

	className, err := s.reportRepo.GetClassName(classID)
	if err != nil {
		return nil, err
	}

	rawRecords, err := s.reportRepo.GetDailyReportRaw(classID, dateStr)
	if err != nil {
		return nil, err
	}

	totalStudents := len(rawRecords)
	summary := domain.DailySummary{}
	records := []domain.DailyRecord{}

	for _, r := range rawRecords {
		records = append(records, domain.DailyRecord{
			NISN:        r.NISN,
			StudentName: r.StudentName,
			Status:      r.Status,
			ScannedAt:   r.ScannedAt,
			IsManual:    r.IsManual,
		})
		switch r.Status {
		case "hadir":
			summary.Hadir++
		case "izin":
			summary.Izin++
		case "sakit":
			summary.Sakit++
		case "tanpa_keterangan":
			summary.TanpaKeterangan++
		}
	}

	if totalStudents > 0 {
		summary.HadirPercentage = float64(summary.Hadir) / float64(totalStudents) * 100
	}

	report := &domain.DailyReport{
		ReportType:    "harian",
		ClassID:       classID,
		ClassName:     className,
		Date:          dateStr,
		TotalStudents: totalStudents,
		Summary:       summary,
		Records:       records,
		GeneratedAt:   time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(report)
	if err == nil {
		s.cache.Set("harian", classID, dateStr, dateStr, generatedBy, data)
	}

	return report, nil
}

func (s *reportService) GetMonthlyReport(classID, monthStr string, forceRefresh bool, generatedBy string) (*domain.MonthlyReport, error) {
	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return nil, errors.New("invalid month format")
	}

	startDate := t.Format("2006-01-02")
	endDate := t.AddDate(0, 1, -1).Format("2006-01-02")

	if forceRefresh {
		s.cache.Delete("bulanan", classID, startDate, endDate)
	} else {
		cached, err := s.cache.Get("bulanan", classID, startDate, endDate)
		if err == nil && cached != nil {
			var report domain.MonthlyReport
			if err := json.Unmarshal(cached, &report); err == nil {
				return &report, nil
			}
		}
	}

	className, err := s.reportRepo.GetClassName(classID)
	if err != nil {
		return nil, err
	}

	rawStats, totalDays, err := s.reportRepo.GetAggregatedReportRaw(classID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	summary := domain.MonthlySummary{
		TotalStudents: len(rawStats),
	}

	var stats []domain.MonthlyStudentStats
	totalHadirPct := 0.0

	for _, r := range rawStats {
		pct := 0.0
		if totalDays > 0 {
			pct = float64(r.Hadir) / float64(totalDays) * 100
		}
		stats = append(stats, domain.MonthlyStudentStats{
			NISN:                 r.NISN,
			StudentName:          r.StudentName,
			Hadir:                r.Hadir,
			Izin:                 r.Izin,
			Sakit:                r.Sakit,
			TanpaKeterangan:      r.TanpaKeterangan,
			DailyStatuses:        r.DailyStatuses,
			AttendancePercentage: pct,
		})

		summary.TotalIzin += r.Izin
		summary.TotalSakit += r.Sakit
		summary.TotalTanpaKeterangan += r.TanpaKeterangan
		totalHadirPct += pct
	}

	if summary.TotalStudents > 0 {
		summary.AvgHadirPercentage = totalHadirPct / float64(summary.TotalStudents)
	}

	monthsId := map[string]string{
		"January": "Januari", "February": "Februari", "March": "Maret",
		"April": "April", "May": "Mei", "June": "Juni",
		"July": "Juli", "August": "Agustus", "September": "September",
		"October": "Oktober", "November": "November", "December": "Desember",
	}
	monthNameEn := t.Format("January")
	monthNameId := monthsId[monthNameEn]
	if monthNameId == "" {
		monthNameId = monthNameEn
	}
	periodName := fmt.Sprintf("%s %d", monthNameId, t.Year())

	report := &domain.MonthlyReport{
		ReportType:   "bulanan",
		ClassID:      classID,
		ClassName:    className,
		Period:       periodName, // e.g., "Mei 2026"
		TotalDays:    totalDays,
		Summary:      summary,
		StudentStats: stats,
		GeneratedAt:  time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(report)
	if err == nil {
		s.cache.Set("bulanan", classID, startDate, endDate, generatedBy, data)
	}

	return report, nil
}

func (s *reportService) GetSemesterReport(classID, academicYear string, semester int, forceRefresh bool, generatedBy string) (*domain.SemesterReport, error) {
	if len(academicYear) != 9 || academicYear[4] != '/' {
		return nil, errors.New("invalid academic year format")
	}

	startYearStr := academicYear[:4]
	endYearStr := academicYear[5:]

	startYear, err1 := strconv.Atoi(startYearStr)
	endYear, err2 := strconv.Atoi(endYearStr)
	if err1 != nil || err2 != nil || startYear+1 != endYear {
		return nil, errors.New("invalid academic year values")
	}

	var startDate, endDate string
	if semester == 1 {
		startDate = fmt.Sprintf("%04d-07-01", startYear)
		endDate = fmt.Sprintf("%04d-12-31", startYear)
	} else if semester == 2 {
		startDate = fmt.Sprintf("%04d-01-01", endYear)
		endDate = fmt.Sprintf("%04d-06-30", endYear)
	} else {
		return nil, errors.New("invalid semester value")
	}

	if forceRefresh {
		s.cache.Delete("semesteran", classID, startDate, endDate)
	} else {
		cached, err := s.cache.Get("semesteran", classID, startDate, endDate)
		if err == nil && cached != nil {
			var report domain.SemesterReport
			if err := json.Unmarshal(cached, &report); err == nil {
				return &report, nil
			}
		}
	}

	className, err := s.reportRepo.GetClassName(classID)
	if err != nil {
		return nil, err
	}

	rawStats, durationDays, err := s.reportRepo.GetAggregatedReportRaw(classID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	rawTrend, err := s.reportRepo.GetTrendRaw(classID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	summary := domain.SemesterSummary{}
	var stats []domain.MonthlyStudentStats
	totalHadirPct := 0.0

	for _, r := range rawStats {
		pct := 0.0
		if durationDays > 0 {
			pct = float64(r.Hadir) / float64(durationDays) * 100
		}
		stats = append(stats, domain.MonthlyStudentStats{
			NISN:                 r.NISN,
			StudentName:          r.StudentName,
			Hadir:                r.Hadir,
			Izin:                 r.Izin,
			Sakit:                r.Sakit,
			TanpaKeterangan:      r.TanpaKeterangan,
			DailyStatuses:        r.DailyStatuses,
			AttendancePercentage: pct,
		})

		summary.TotalIzin += r.Izin
		summary.TotalSakit += r.Sakit
		summary.TotalTanpaKeterangan += r.TanpaKeterangan
		totalHadirPct += pct
	}

	if len(stats) > 0 {
		summary.AvgAttendance = totalHadirPct / float64(len(stats))
	}

	var trend []domain.SemesterTrend
	for _, tr := range rawTrend {
		pct := 0.0
		if tr.Total > 0 {
			pct = float64(tr.Hadir) / float64(tr.Total) * 100
		}
		trend = append(trend, domain.SemesterTrend{
			Month:                tr.Month,
			AttendancePercentage: pct,
		})
	}

	report := &domain.SemesterReport{
		ReportType:   "semesteran",
		ClassID:      classID,
		ClassName:    className,
		Period:       fmt.Sprintf("Semester %d - %s", semester, academicYear),
		DurationDays: durationDays,
		Summary:      summary,
		Trend:        trend,
		StudentStats: stats,
		GeneratedAt:  time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(report)
	if err == nil {
		s.cache.Set("semesteran", classID, startDate, endDate, generatedBy, data)
	}

	return report, nil
}

func (s *reportService) GetStudentReport(studentID, startDate, endDate string) (*domain.StudentReport, error) {
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	records, err := s.reportRepo.GetStudentReportRaw(studentID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	summary := domain.StudentSummary{}
	var studentName, nisn, className string

	for _, r := range records {
		if studentName == "" {
			studentName = r.StudentName
			nisn = r.NISN
			className = r.ClassName
		}
		switch r.Status {
		case "hadir":
			summary.Hadir++
		case "izin":
			summary.Izin++
		case "sakit":
			summary.Sakit++
		case "tanpa_keterangan":
			summary.TanpaKeterangan++
		}
	}

	total := summary.Hadir + summary.Izin + summary.Sakit + summary.TanpaKeterangan
	if total > 0 {
		summary.HadirPercentage = float64(summary.Hadir) / float64(total) * 100
	}

	return &domain.StudentReport{
		StudentID:   studentID,
		StudentName: studentName,
		NISN:        nisn,
		ClassName:   className,
		Summary:     summary,
		Records:     records,
	}, nil
}
