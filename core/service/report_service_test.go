package service_test

import (
	"encoding/json"
	"testing"

	"alhikmah-attendance-api/core/domain"
	"alhikmah-attendance-api/core/mocks"
	"alhikmah-attendance-api/core/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReportService_GetDailyReport_ForceRefresh(t *testing.T) {
	mockRepo := new(mocks.ReportRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockCache := new(mocks.ReportCacheRepository)

	svc := service.NewReportService(mockRepo, mockStudentRepo, mockCache)

	dateStr := "2026-05-30"
	classID := "class1"
	
	mockCache.On("Delete", "harian", classID, dateStr, dateStr).Return(nil)
	mockRepo.On("GetClassName", classID).Return("10A", nil)
	mockRepo.On("GetDailyReportRaw", classID, dateStr).Return([]domain.DailyRecordRaw{
		{NISN: "123", StudentName: "John", Status: "hadir"},
	}, nil)
	mockCache.On("Set", "harian", classID, dateStr, dateStr, "admin", mock.Anything).Return(nil)

	report, err := svc.GetDailyReport(classID, dateStr, true, "admin")
	
	assert.NoError(t, err)
	assert.Equal(t, "10A", report.ClassName)
	assert.Equal(t, 1, report.TotalStudents)
	assert.Equal(t, 1, report.Summary.Hadir)
	assert.Equal(t, 100.0, report.Summary.HadirPercentage)
}

func TestReportService_GetDailyReport_Cached(t *testing.T) {
	mockRepo := new(mocks.ReportRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockCache := new(mocks.ReportCacheRepository)

	svc := service.NewReportService(mockRepo, mockStudentRepo, mockCache)

	dateStr := "2026-05-30"
	classID := "class1"
	
	cachedReport := domain.DailyReport{ClassName: "10A Cached"}
	cachedData, _ := json.Marshal(cachedReport)
	
	mockCache.On("Get", "harian", classID, dateStr, dateStr).Return(json.RawMessage(cachedData), nil)

	report, err := svc.GetDailyReport(classID, dateStr, false, "admin")
	
	assert.NoError(t, err)
	assert.Equal(t, "10A Cached", report.ClassName)
}

func TestReportService_GetMonthlyReport(t *testing.T) {
	mockRepo := new(mocks.ReportRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockCache := new(mocks.ReportCacheRepository)

	svc := service.NewReportService(mockRepo, mockStudentRepo, mockCache)

	classID := "class1"
	monthStr := "2026-05"
	startDate := "2026-05-01"
	endDate := "2026-05-31"

	mockCache.On("Delete", "bulanan", classID, startDate, endDate).Return(nil)
	mockRepo.On("GetClassName", classID).Return("10A", nil)
	
	rawStats := []domain.MonthlyStatRaw{
		{NISN: "123", StudentName: "John", Hadir: 20},
	}
	mockRepo.On("GetAggregatedReportRaw", classID, startDate, endDate).Return(rawStats, 20, nil)
	mockCache.On("Set", "bulanan", classID, startDate, endDate, "admin", mock.Anything).Return(nil)

	report, err := svc.GetMonthlyReport(classID, monthStr, true, "admin")
	
	assert.NoError(t, err)
	assert.Equal(t, "10A", report.ClassName)
	assert.Equal(t, 20, report.TotalDays)
	assert.Equal(t, 1, report.Summary.TotalStudents)
	assert.Equal(t, 100.0, report.StudentStats[0].AttendancePercentage)
}

func TestReportService_GetSemesterReport(t *testing.T) {
	mockRepo := new(mocks.ReportRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockCache := new(mocks.ReportCacheRepository)

	svc := service.NewReportService(mockRepo, mockStudentRepo, mockCache)

	classID := "class1"
	academicYear := "2025/2026"
	semester := 2
	startDate := "2026-01-01"
	endDate := "2026-06-30"

	mockCache.On("Delete", "semesteran", classID, startDate, endDate).Return(nil)
	mockRepo.On("GetClassName", classID).Return("10A", nil)
	
	mockRepo.On("GetAggregatedReportRaw", classID, startDate, endDate).Return([]domain.MonthlyStatRaw{}, 100, nil)
	mockRepo.On("GetTrendRaw", classID, startDate, endDate).Return([]domain.TrendRaw{}, nil)
	mockCache.On("Set", "semesteran", classID, startDate, endDate, "admin", mock.Anything).Return(nil)

	report, err := svc.GetSemesterReport(classID, academicYear, semester, true, "admin")
	
	assert.NoError(t, err)
	assert.Equal(t, "10A", report.ClassName)
}

func TestReportService_GetSemesterReport_InvalidYear(t *testing.T) {
	svc := service.NewReportService(nil, nil, nil)
	_, err := svc.GetSemesterReport("1", "20252026", 1, true, "admin")
	assert.Error(t, err)
}

func TestReportService_GetStudentReport(t *testing.T) {
	mockRepo := new(mocks.ReportRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	mockCache := new(mocks.ReportCacheRepository)

	svc := service.NewReportService(mockRepo, mockStudentRepo, mockCache)

	records := []domain.StudentReportRecord{
		{StudentName: "John", NISN: "123", ClassName: "10A", Status: "hadir"},
	}
	mockRepo.On("GetStudentReportRaw", "s1", "2026-01-01", "2026-05-30").Return(records, nil)

	report, err := svc.GetStudentReport("s1", "2026-01-01", "2026-05-30")
	
	assert.NoError(t, err)
	assert.Equal(t, "John", report.StudentName)
	assert.Equal(t, 100.0, report.Summary.HadirPercentage)
}
