package repository

import (
	"database/sql"
	"time"

	"alhikmah-attendance-api/internal/domain"
)

type reportPostgres struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) domain.ReportRepository {
	return &reportPostgres{db: db}
}

func (r *reportPostgres) GetClassName(classID string) (string, error) {
	var name string
	err := r.db.QueryRow("SELECT class_name FROM classes WHERE id = $1", classID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "Unknown Class", nil
		}
		return "", err
	}
	return name, nil
}

func (r *reportPostgres) GetDailyReportRaw(classID, dateStr string) ([]domain.DailyRecordRaw, error) {
	query := `
		SELECT s.nisn, s.full_name, COALESCE(CAST(a.status AS TEXT), 'belum_absen') as status, a.scanned_at, COALESCE(a.is_manual, false)
		FROM students s
		LEFT JOIN attendances a ON s.id = a.student_id AND a.attendance_date = $1
		WHERE s.class_id = $2 AND s.is_active = true
		ORDER BY s.full_name ASC
	`
	rows, err := r.db.Query(query, dateStr, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.DailyRecordRaw
	for rows.Next() {
		var item domain.DailyRecordRaw
		var scannedAt sql.NullTime
		if err := rows.Scan(&item.NISN, &item.StudentName, &item.Status, &scannedAt, &item.IsManual); err != nil {
			return nil, err
		}
		if scannedAt.Valid {
			t := scannedAt.Time.Format(time.RFC3339)
			item.ScannedAt = &t
		}
		records = append(records, item)
	}

	return records, nil
}

func (r *reportPostgres) GetAggregatedReportRaw(classID, startDate, endDate string) ([]domain.MonthlyStatRaw, int, error) {
	// First get the student stats
	query := `
		SELECT s.nisn, s.full_name,
			COUNT(CASE WHEN a.status = 'hadir' THEN 1 END) as total_hadir,
			COUNT(CASE WHEN a.status = 'izin' THEN 1 END) as total_izin,
			COUNT(CASE WHEN a.status = 'sakit' THEN 1 END) as total_sakit,
			COUNT(CASE WHEN a.status = 'tanpa_keterangan' THEN 1 END) as total_tanpa_keterangan
		FROM students s
		LEFT JOIN attendances a ON s.id = a.student_id AND a.attendance_date >= $1 AND a.attendance_date <= $2
		WHERE s.class_id = $3 AND s.is_active = true
		GROUP BY s.id, s.nisn, s.full_name
		ORDER BY s.full_name ASC
	`
	rows, err := r.db.Query(query, startDate, endDate, classID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stats []domain.MonthlyStatRaw
	for rows.Next() {
		var item domain.MonthlyStatRaw
		if err := rows.Scan(&item.NISN, &item.StudentName, &item.Hadir, &item.Izin, &item.Sakit, &item.TanpaKeterangan); err != nil {
			return nil, 0, err
		}
		stats = append(stats, item)
	}

	// Calculate total working days in this period (days where at least 1 attendance exists in this class)
	var totalDays int
	queryDays := `
		SELECT COUNT(DISTINCT attendance_date)
		FROM attendances
		WHERE class_id = $1 AND attendance_date >= $2 AND attendance_date <= $3
	`
	err = r.db.QueryRow(queryDays, classID, startDate, endDate).Scan(&totalDays)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return stats, totalDays, nil
}

func (r *reportPostgres) GetTrendRaw(classID, startDate, endDate string) ([]domain.TrendRaw, error) {
	query := `
		SELECT TO_CHAR(attendance_date, 'YYYY-MM') as month,
			   COUNT(CASE WHEN status = 'hadir' THEN 1 END) as hadir,
			   COUNT(*) as total
		FROM attendances
		WHERE class_id = $1 AND attendance_date >= $2 AND attendance_date <= $3
		GROUP BY TO_CHAR(attendance_date, 'YYYY-MM')
		ORDER BY month ASC
	`
	rows, err := r.db.Query(query, classID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []domain.TrendRaw
	for rows.Next() {
		var item domain.TrendRaw
		if err := rows.Scan(&item.Month, &item.Hadir, &item.Total); err != nil {
			return nil, err
		}
		trends = append(trends, item)
	}
	return trends, nil
}
