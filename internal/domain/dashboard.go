package domain

import "time"

type RecentActivity struct {
	ID          string    `json:"id"`
	StudentName string    `json:"student_name"`
	ClassName   string    `json:"class_name"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

type AttendanceTrend struct {
	Date            string `json:"date"`
	Hadir           int    `json:"hadir"`
	Izin            int    `json:"izin"`
	Sakit           int    `json:"sakit"`
	TanpaKeterangan int    `json:"tanpa_keterangan"`
}

type DashboardRepository interface {
	GetRecentActivity(limit int) ([]RecentActivity, error)
	GetAttendanceTrend(startDate, endDate time.Time) ([]AttendanceTrend, error)
}

type DashboardService interface {
	GetRecentActivity(limit int) ([]RecentActivity, error)
	GetAttendanceTrend(days int) ([]AttendanceTrend, error)
}
