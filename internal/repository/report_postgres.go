package repository

import (
	"database/sql"
	"encoding/json"
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
	err := r.db.QueryRow("SELECT CONCAT('Kelas ', grade, ' ', room_name, COALESCE(' ' || section, '')) FROM classes WHERE id = $1", classID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "Unknown Class", nil
		}
		return "", err
	}
	return name, nil
}

func (r *reportPostgres) GetDailyReportRaw(classID, dateStr, subject string) ([]domain.DailyRecordRaw, error) {
	query := `
		SELECT s.nisn, s.full_name, COALESCE(CAST(a.status AS TEXT), 'belum_absen') as status, a.scanned_at, COALESCE(a.is_manual, false)
		FROM students s
		LEFT JOIN attendances a ON s.id = a.student_id AND a.attendance_date = $1
	`
	if subject == "" {
		query += " AND a.subject IS NULL"
	} else {
		query += " AND a.subject = '" + subject + "'" // Safe since subject comes from query param, but better to parameterize. Let's do parameterization.
	}
	
	// Better parameterization
	query = `
		SELECT s.nisn, s.full_name, COALESCE(CAST(a.status AS TEXT), 'belum_absen') as status, a.scanned_at, COALESCE(a.is_manual, false)
		FROM student_enrollments e
		JOIN students s ON e.student_id = s.id
		LEFT JOIN attendances a ON s.id = a.student_id AND a.attendance_date = $1 AND a.class_id = $2
	`
	if subject == "" {
		query += " AND a.subject IS NULL"
	} else {
		query += " AND a.subject = $3"
	}
	query += ` WHERE e.class_id = $2 AND s.is_active = true
		AND e.enrolled_at::DATE <= $1 
		AND (e.ended_at IS NULL OR e.ended_at::DATE >= $1)
		ORDER BY s.full_name ASC`
		
	var rows *sql.Rows
	var err error
	if subject == "" {
		rows, err = r.db.Query(query, dateStr, classID)
	} else {
		rows, err = r.db.Query(query, dateStr, classID, subject)
	}
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

func (r *reportPostgres) GetAggregatedReportRaw(classID, startDate, endDate, subject string) ([]domain.MonthlyStatRaw, int, error) {
	// First get the student stats
	query := `
		SELECT s.nisn, s.full_name,
			COUNT(CASE WHEN a.status = 'hadir' THEN 1 END) as total_hadir,
			COUNT(CASE WHEN a.status = 'izin' THEN 1 END) as total_izin,
			COUNT(CASE WHEN a.status = 'sakit' THEN 1 END) as total_sakit,
			COUNT(CASE WHEN a.status = 'tanpa_keterangan' THEN 1 END) as total_tanpa_keterangan,
			COALESCE(
				JSON_OBJECT_AGG(
					EXTRACT(DAY FROM a.attendance_date)::int, 
					a.status
				) FILTER (WHERE a.attendance_date IS NOT NULL), 
				'{}'
			) as daily_statuses
		FROM student_enrollments e
		JOIN students s ON e.student_id = s.id
		LEFT JOIN attendances a ON s.id = a.student_id AND a.class_id = $3 AND a.attendance_date >= $1 AND a.attendance_date <= $2
	`
	if subject == "" {
		query += " AND a.subject IS NULL"
	} else {
		query += " AND a.subject = $4"
	}
	query += ` WHERE e.class_id = $3 AND s.is_active = true
		AND e.enrolled_at::DATE <= $2
		AND (e.ended_at IS NULL OR e.ended_at::DATE >= $1)
		GROUP BY s.id, s.nisn, s.full_name
		ORDER BY s.full_name ASC
	`
	var rows *sql.Rows
	var err error
	if subject == "" {
		rows, err = r.db.Query(query, startDate, endDate, classID)
	} else {
		rows, err = r.db.Query(query, startDate, endDate, classID, subject)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stats []domain.MonthlyStatRaw
	for rows.Next() {
		var item domain.MonthlyStatRaw
		var dailyStatusesJSON []byte
		if err := rows.Scan(
			&item.NISN, &item.StudentName, 
			&item.Hadir, &item.Izin, &item.Sakit, &item.TanpaKeterangan,
			&dailyStatusesJSON,
		); err != nil {
			return nil, 0, err
		}
		
		item.DailyStatuses = make(map[int]string)
		if len(dailyStatusesJSON) > 0 && string(dailyStatusesJSON) != "{}" {
			importJsonErr := json.Unmarshal(dailyStatusesJSON, &item.DailyStatuses)
			if importJsonErr != nil {
				// Ignore JSON unmarshal error, it will just leave it empty
			}
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
	if subject == "" {
		queryDays += " AND subject IS NULL"
		err = r.db.QueryRow(queryDays, classID, startDate, endDate).Scan(&totalDays)
	} else {
		queryDays += " AND subject = $4"
		err = r.db.QueryRow(queryDays, classID, startDate, endDate, subject).Scan(&totalDays)
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return stats, totalDays, nil
}

func (r *reportPostgres) GetTrendRaw(classID, startDate, endDate, subject string) ([]domain.TrendRaw, error) {
	query := `
		SELECT TO_CHAR(attendance_date, 'YYYY-MM') as month,
			   COUNT(CASE WHEN status = 'hadir' THEN 1 END) as hadir,
			   COUNT(*) as total
		FROM attendances
		WHERE class_id = $1 AND attendance_date >= $2 AND attendance_date <= $3
	`
	if subject == "" {
		query += " AND subject IS NULL"
	} else {
		query += " AND subject = $4"
	}
	query += ` GROUP BY TO_CHAR(attendance_date, 'YYYY-MM')
		ORDER BY month ASC
	`
	
	var rows *sql.Rows
	var err error
	if subject == "" {
		rows, err = r.db.Query(query, classID, startDate, endDate)
	} else {
		rows, err = r.db.Query(query, classID, startDate, endDate, subject)
	}
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
func (r *reportPostgres) GetStudentReportRaw(studentID, startDate, endDate string) ([]domain.StudentReportRecordRaw, error) {
	query := `
		SELECT 
			a.id, 
			s.id as student_id, 
			s.full_name as student_name, 
			s.nisn, 
			c.id as class_id, 
			CONCAT('Kelas ', c.grade, ' ', c.room_name, COALESCE(' ' || c.section, '')) as class_name, 
			TO_CHAR(a.attendance_date, 'YYYY-MM-DD') as attendance_date, 
			a.status, 
			a.recorded_by, 
			TO_CHAR(a.recorded_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as recorded_at, 
			a.scanned_at, 
			a.notes, 
			a.is_manual
		FROM attendances a
		JOIN students s ON a.student_id = s.id
		JOIN classes c ON a.class_id = c.id
		WHERE a.student_id = $1 AND a.attendance_date >= $2 AND a.attendance_date <= $3
		ORDER BY a.attendance_date DESC
	`
	rows, err := r.db.Query(query, studentID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.StudentReportRecordRaw
	for rows.Next() {
		var item domain.StudentReportRecordRaw
		var scannedAt sql.NullTime
		var notes sql.NullString
		
		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.StudentName,
			&item.NISN,
			&item.ClassID,
			&item.ClassName,
			&item.AttendanceDate,
			&item.Status,
			&item.RecordedBy,
			&item.RecordedAt,
			&scannedAt,
			&notes,
			&item.IsManual,
		); err != nil {
			return nil, err
		}

		if scannedAt.Valid {
			t := scannedAt.Time.Format(time.RFC3339)
			item.ScannedAt = &t
		}
		if notes.Valid {
			n := notes.String
			item.Notes = &n
		}

		records = append(records, item)
	}

	return records, nil
}
