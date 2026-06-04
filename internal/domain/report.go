package domain

import (
        "encoding/json"
)

// Raw Internal structs for repository
type DailyRecordRaw struct {
	NISN        string
	StudentName string
	Status      string
	ScannedAt   *string
	IsManual    bool
}

type MonthlyStatRaw struct {
	NISN            string
	StudentName     string
	Hadir           int
	Izin            int
	Sakit           int
	TanpaKeterangan int
	DailyStatuses   map[int]string
}

type TrendRaw struct {
	Month string
	Hadir int
	Total int
}

type StudentReportRecordRaw struct {
	ID             string  `json:"id"`
	StudentID      string  `json:"student_id"`
	StudentName    string  `json:"student_name"`
	NISN           string  `json:"nisn"`
	ClassID        string  `json:"class_id"`
	ClassName      string  `json:"class_name"`
	AttendanceDate string  `json:"attendance_date"`
	Status         string  `json:"status"`
	RecordedBy     string  `json:"recorded_by"`
	RecordedAt     string  `json:"recorded_at"`
	ScannedAt      *string `json:"scanned_at"`
	Notes          *string `json:"notes"`
	IsManual       bool    `json:"is_manual"`
}

type DailySummary struct {
	Hadir           int
	Izin            int
	Sakit           int
	TanpaKeterangan int
	HadirPercentage float64
}

type DailyRecord struct {
	NISN        string
	StudentName string
	Status      string
	ScannedAt   *string
	IsManual    bool
}

type DailyReport struct {
	ReportType    string
	ClassID       string
	ClassName     string
	Date          string
	TotalStudents int
	Summary       DailySummary
	Records       []DailyRecord
	Subject       *string
	GeneratedAt   string
}

type MonthlyStudentStats struct {
	NISN                 string
	StudentName          string
	Hadir                int
	Izin                 int
	Sakit                int
	TanpaKeterangan      int
	DailyStatuses        map[int]string
	AttendancePercentage float64
}

type MonthlySummary struct {
	TotalStudents        int
	AvgHadirPercentage   float64
	TotalIzin            int
	TotalSakit           int
	TotalTanpaKeterangan int
}

type MonthlyReport struct {
	ReportType   string
	ClassID      string
	ClassName    string
	Period       string
	TotalDays    int
	Summary      MonthlySummary
	StudentStats []MonthlyStudentStats
	Subject      *string
	GeneratedAt  string
}

type SemesterTrend struct {
	Month                string
	AttendancePercentage float64
}

type SemesterSummary struct {
	AvgAttendance        float64
	TotalIzin            int
	TotalSakit           int
	TotalTanpaKeterangan int
}

type SemesterReport struct {
	ReportType   string
	ClassID      string
	ClassName    string
	Period       string
	DurationDays int
	Summary      SemesterSummary
	Trend        []SemesterTrend
	StudentStats []MonthlyStudentStats
	Subject      *string
	GeneratedAt  string
}

type StudentSummary struct {
	Hadir           int
	Izin            int
	Sakit           int
	TanpaKeterangan int
	HadirPercentage float64
}

type StudentReportRecord struct {
	ID             string
	StudentID      string
	StudentName    string
	NISN           string
	ClassID        string
	ClassName      string
	AttendanceDate string
	Status         string
	RecordedBy     string
	RecordedAt     string
	ScannedAt      *string
	Notes          *string
	IsManual       bool
}

type StudentReport struct {
	StudentID   string
	StudentName string
	NISN        string
	ClassName   string
	Summary     StudentSummary
	Records     []StudentReportRecord
}


type ReportRepository interface {
	GetClassName(classID string) (string, error)
	GetDailyReportRaw(classID, dateStr, subject string) ([]DailyRecordRaw, error)
	GetAggregatedReportRaw(classID, startDate, endDate, subject string) ([]MonthlyStatRaw, int, error)
	GetTrendRaw(classID, startDate, endDate, subject string) ([]TrendRaw, error)
	GetStudentReportRaw(studentID, startDate, endDate string) ([]StudentReportRecordRaw, error)
}

type ReportCacheRepository interface {
	Get(reportType, classID, periodStart, periodEnd, subject string) (json.RawMessage, error)
	Set(reportType, classID, periodStart, periodEnd, subject, generatedBy string, data json.RawMessage) error
	Delete(reportType, classID, periodStart, periodEnd, subject string) error
}

type ReportService interface {
	GetDailyReport(classID, dateStr, subject string, forceRefresh bool, generatedBy string) (*DailyReport, error)
	GetMonthlyReport(classID, monthStr, subject string, forceRefresh bool, generatedBy string) (*MonthlyReport, error)
	GetSemesterReport(classID, academicYear string, semester int, subject string, forceRefresh bool, generatedBy string) (*SemesterReport, error)
	GetStudentReport(studentID, startDate, endDate string) (*StudentReport, error)
}
