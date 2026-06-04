package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToDailyReportDTO(r *domain.DailyReport) *dto.DailyReportResponse {
	records := make([]dto.DailyRecord, len(r.Records))
	for i, v := range r.Records {
		records[i] = dto.DailyRecord{
			NISN:        v.NISN,
			StudentName: v.StudentName,
			Status:      v.Status,
			ScannedAt:   v.ScannedAt,
			IsManual:    v.IsManual,
		}
	}

	return &dto.DailyReportResponse{
		ReportType:    r.ReportType,
		ClassID:       r.ClassID,
		ClassName:     r.ClassName,
		Date:          r.Date,
		TotalStudents: r.TotalStudents,
		Summary: dto.DailySummary{
			Hadir:           r.Summary.Hadir,
			Izin:            r.Summary.Izin,
			Sakit:           r.Summary.Sakit,
			TanpaKeterangan: r.Summary.TanpaKeterangan,
			HadirPercentage: r.Summary.HadirPercentage,
		},
		Records:     records,
		Subject:     r.Subject,
		GeneratedAt: r.GeneratedAt,
	}
}

func ToMonthlyReportDTO(r *domain.MonthlyReport) *dto.MonthlyReportResponse {
	stats := make([]dto.MonthlyStudentStats, len(r.StudentStats))
	for i, v := range r.StudentStats {
		stats[i] = dto.MonthlyStudentStats{
			NISN:                 v.NISN,
			StudentName:          v.StudentName,
			Hadir:                v.Hadir,
			Izin:                 v.Izin,
			Sakit:                v.Sakit,
			TanpaKeterangan:      v.TanpaKeterangan,
			DailyStatuses:        v.DailyStatuses,
			AttendancePercentage: v.AttendancePercentage,
		}
	}

	return &dto.MonthlyReportResponse{
		ReportType: r.ReportType,
		ClassID:    r.ClassID,
		ClassName:  r.ClassName,
		Period:     r.Period,
		TotalDays:  r.TotalDays,
		Summary: dto.MonthlySummary{
			TotalStudents:        r.Summary.TotalStudents,
			AvgHadirPercentage:   r.Summary.AvgHadirPercentage,
			TotalIzin:            r.Summary.TotalIzin,
			TotalSakit:           r.Summary.TotalSakit,
			TotalTanpaKeterangan: r.Summary.TotalTanpaKeterangan,
		},
		StudentStats: stats,
		Subject:      r.Subject,
		GeneratedAt:  r.GeneratedAt,
	}
}

func ToSemesterReportDTO(r *domain.SemesterReport) *dto.SemesterReportResponse {
	stats := make([]dto.MonthlyStudentStats, len(r.StudentStats))
	for i, v := range r.StudentStats {
		stats[i] = dto.MonthlyStudentStats{
			NISN:                 v.NISN,
			StudentName:          v.StudentName,
			Hadir:                v.Hadir,
			Izin:                 v.Izin,
			Sakit:                v.Sakit,
			TanpaKeterangan:      v.TanpaKeterangan,
			DailyStatuses:        v.DailyStatuses,
			AttendancePercentage: v.AttendancePercentage,
		}
	}

	trend := make([]dto.SemesterTrend, len(r.Trend))
	for i, v := range r.Trend {
		trend[i] = dto.SemesterTrend{
			Month:                v.Month,
			AttendancePercentage: v.AttendancePercentage,
		}
	}

	return &dto.SemesterReportResponse{
		ReportType:   r.ReportType,
		ClassID:      r.ClassID,
		ClassName:    r.ClassName,
		Period:       r.Period,
		DurationDays: r.DurationDays,
		Summary: dto.SemesterSummary{
			AvgAttendance:        r.Summary.AvgAttendance,
			TotalIzin:            r.Summary.TotalIzin,
			TotalSakit:           r.Summary.TotalSakit,
			TotalTanpaKeterangan: r.Summary.TotalTanpaKeterangan,
		},
		Trend:        trend,
		StudentStats: stats,
		Subject:      r.Subject,
		GeneratedAt:  r.GeneratedAt,
	}
}

func ToStudentReportDTO(r *domain.StudentReport) *dto.StudentReportResponse {
	records := make([]dto.StudentReportRecord, len(r.Records))
	for i, v := range r.Records {
		records[i] = dto.StudentReportRecord{
			ID:             v.ID,
			StudentID:      v.StudentID,
			StudentName:    v.StudentName,
			NISN:           v.NISN,
			ClassID:        v.ClassID,
			ClassName:      v.ClassName,
			AttendanceDate: v.AttendanceDate,
			Status:         v.Status,
			RecordedBy:     v.RecordedBy,
			RecordedAt:     v.RecordedAt,
			ScannedAt:      v.ScannedAt,
			Notes:          v.Notes,
			IsManual:       v.IsManual,
		}
	}

	return &dto.StudentReportResponse{
		StudentID:   r.StudentID,
		StudentName: r.StudentName,
		NISN:        r.NISN,
		ClassName:   r.ClassName,
		Summary: dto.StudentSummary{
			Hadir:           r.Summary.Hadir,
			Izin:            r.Summary.Izin,
			Sakit:           r.Summary.Sakit,
			TanpaKeterangan: r.Summary.TanpaKeterangan,
			HadirPercentage: r.Summary.HadirPercentage,
		},
		Records: records,
	}
}
