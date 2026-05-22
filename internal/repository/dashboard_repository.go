package repository

import (
	"database/sql"
	"time"

	"alhikmah-attendance-api/internal/domain"
)

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetRecentActivity(limit int) ([]domain.RecentActivity, error) {
	query := `
		SELECT 
			a.id, 
			s.full_name, 
			c.class_name, 
			a.status, 
			a.date 
		FROM attendances a
		JOIN students s ON a.student_id = s.id
		JOIN classes c ON s.class_id = c.id
		ORDER BY a.updated_at DESC, a.created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []domain.RecentActivity
	for rows.Next() {
		var act domain.RecentActivity
		if err := rows.Scan(&act.ID, &act.StudentName, &act.ClassName, &act.Status, &act.Timestamp); err != nil {
			return nil, err
		}
		activities = append(activities, act)
	}

	return activities, nil
}

func (r *dashboardRepository) GetAttendanceTrend(startDate, endDate time.Time) ([]domain.AttendanceTrend, error) {
	query := `
		SELECT 
			TO_CHAR(date, 'YYYY-MM-DD') as fmt_date,
			COUNT(CASE WHEN status = 'hadir' THEN 1 END) as hadir,
			COUNT(CASE WHEN status = 'izin' THEN 1 END) as izin,
			COUNT(CASE WHEN status = 'sakit' THEN 1 END) as sakit,
			COUNT(CASE WHEN status = 'tanpa_keterangan' THEN 1 END) as tanpa_keterangan
		FROM attendances
		WHERE date >= $1 AND date <= $2
		GROUP BY date
		ORDER BY date ASC
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []domain.AttendanceTrend
	for rows.Next() {
		var t domain.AttendanceTrend
		if err := rows.Scan(&t.Date, &t.Hadir, &t.Izin, &t.Sakit, &t.TanpaKeterangan); err != nil {
			return nil, err
		}
		trends = append(trends, t)
	}

	return trends, nil
}
