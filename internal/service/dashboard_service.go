package service

import (
	"time"

	"alhikmah-attendance-api/internal/domain"
)

type dashboardService struct {
	repo domain.DashboardRepository
}

func NewDashboardService(repo domain.DashboardRepository) domain.DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetRecentActivity(limit int) ([]domain.RecentActivity, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetRecentActivity(limit)
}

func (s *dashboardService) GetAttendanceTrend(days int) ([]domain.AttendanceTrend, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.repo.GetAttendanceTrend(startDate, endDate)
}
