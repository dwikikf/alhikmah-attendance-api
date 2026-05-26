package domain

import "encoding/json"

type DailyReport struct {
	ReportType    string         `json:"report_type"` // "harian"
	ClassID       string         `json:"class_id"`
	ClassName     string         `json:"class_name"`
	Date          string         `json:"date"`
	TotalStudents int            `json:"total_students"`
	Summary       DailySummary   `json:"summary"`
	Records       []DailyRecord  `json:"records"`
	GeneratedAt   string         `json:"generated_at"`
}

type DailySummary struct {
	Hadir           int     `json:"hadir"`
	Izin            int     `json:"izin"`
	Sakit           int     `json:"sakit"`
	TanpaKeterangan int     `json:"tanpa_keterangan"`
	HadirPercentage float64 `json:"hadir_percentage"`
}

type DailyRecord struct {
	NISN        string  `json:"nisn"`
	StudentName string  `json:"student_name"`
	Status      string  `json:"status"`
	ScannedAt   *string `json:"scanned_at"`
	IsManual    bool    `json:"is_manual"`
}

type MonthlyStudentStats struct {
	NISN                 string  `json:"nisn"`
	StudentName          string  `json:"student_name"`
	Hadir                int     `json:"hadir"`
	Izin                 int     `json:"izin"`
	Sakit                int            `json:"sakit"`
	TanpaKeterangan      int            `json:"tanpa_keterangan"`
	DailyStatuses        map[int]string `json:"daily_statuses"`
	AttendancePercentage float64        `json:"attendance_percentage"`
}

type MonthlySummary struct {
	TotalStudents        int     `json:"total_students"`
	AvgHadirPercentage   float64 `json:"avg_hadir_percentage"`
	TotalIzin            int     `json:"total_izin"`
	TotalSakit           int     `json:"total_sakit"`
	TotalTanpaKeterangan int     `json:"total_tanpa_keterangan"`
}

type MonthlyReport struct {
	ReportType   string                `json:"report_type"` // "bulanan"
	ClassID      string                `json:"class_id"`
	ClassName    string                `json:"class_name"`
	Period       string                `json:"period"` // e.g., "May 2024"
	TotalDays    int                   `json:"total_days"`
	Summary      MonthlySummary        `json:"summary"`
	StudentStats []MonthlyStudentStats `json:"student_stats"`
	GeneratedAt  string                `json:"generated_at"`
}

type SemesterTrend struct {
	Month                string  `json:"month"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

type SemesterSummary struct {
	AvgAttendance        float64 `json:"avg_attendance"`
	TotalIzin            int     `json:"total_izin"`
	TotalSakit           int     `json:"total_sakit"`
	TotalTanpaKeterangan int     `json:"total_tanpa_keterangan"`
}

type SemesterReport struct {
	ReportType   string                `json:"report_type"` // "semesteran"
	ClassID      string                `json:"class_id"`
	ClassName    string                `json:"class_name"`
	Period       string                `json:"period"` // e.g., "Semester 1 - 2024/2025"
	DurationDays int                   `json:"duration_days"`
	Summary      SemesterSummary       `json:"summary"`
	Trend        []SemesterTrend       `json:"trend"`
	StudentStats []MonthlyStudentStats `json:"student_stats"`
	GeneratedAt  string                `json:"generated_at"`
}

// Internal structs for repository
type StudentSummary struct {
	Hadir           int     `json:"hadir"`
	Izin            int     `json:"izin"`
	Sakit           int     `json:"sakit"`
	TanpaKeterangan int     `json:"tanpa_keterangan"`
	HadirPercentage float64 `json:"hadir_percentage"`
}

type StudentReportRecord struct {
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

type StudentReport struct {
	StudentID   string                `json:"student_id"`
	StudentName string                `json:"student_name"`
	NISN        string                `json:"nisn"`
	ClassName   string                `json:"class_name"`
	Summary     StudentSummary        `json:"summary"`
	Records     []StudentReportRecord `json:"records"`
}

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

type ReportRepository interface {
	GetClassName(classID string) (string, error)
	GetDailyReportRaw(classID, dateStr string) ([]DailyRecordRaw, error)
	GetAggregatedReportRaw(classID, startDate, endDate string) ([]MonthlyStatRaw, int, error)
	GetTrendRaw(classID, startDate, endDate string) ([]TrendRaw, error)
	GetStudentReportRaw(studentID, startDate, endDate string) ([]StudentReportRecord, error)
}

type ReportCacheRepository interface {
	Get(reportType, classID, periodStart, periodEnd string) (json.RawMessage, error)
	Set(reportType, classID, periodStart, periodEnd, generatedBy string, data json.RawMessage) error
}


type ReportService interface {
	GetDailyReport(classID, dateStr string) (*DailyReport, error)
	GetMonthlyReport(classID, monthStr string) (*MonthlyReport, error)
	GetSemesterReport(classID, academicYear string, semester int) (*SemesterReport, error)
	GetStudentReport(studentID, startDate, endDate string) (*StudentReport, error)
}
