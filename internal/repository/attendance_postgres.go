package repository

import (
	"database/sql"
	"time"

	"alhikmah-attendance-api/internal/domain"
)

type attendancePostgres struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) domain.AttendanceRepository {
	return &attendancePostgres{db: db}
}

func (r *attendancePostgres) MarkAttendance(attendance *domain.Attendance) error {
	query := `
		INSERT INTO attendances (student_id, class_id, attendance_date, status, recorded_by, scanned_at, is_manual, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, recorded_at
	`
	return r.db.QueryRow(
		query,
		attendance.StudentID,
		attendance.ClassID,
		attendance.AttendanceDate,
		attendance.Status,
		attendance.RecordedBy,
		attendance.ScannedAt,
		attendance.IsManual,
		attendance.Notes,
	).Scan(&attendance.ID, &attendance.RecordedAt)
}

func (r *attendancePostgres) UpdateAttendance(attendance *domain.Attendance, audit *domain.AttendanceAudit) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update attendance
	updateQuery := `
		UPDATE attendances 
		SET status = $1, recorded_by = $2, is_manual = true
		WHERE id = $3
		RETURNING status
	`
	var oldStatus string
	err = tx.QueryRow(updateQuery, attendance.Status, attendance.RecordedBy, attendance.ID).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// Insert audit
	auditQuery := `
		INSERT INTO attendance_audits (attendance_id, old_status, new_status, changed_by, reason)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(auditQuery, audit.AttendanceID, audit.OldStatus, audit.NewStatus, audit.ChangedBy, audit.Reason)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *attendancePostgres) GetByClassAndDate(classID string, date time.Time) ([]*domain.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, attendance_date, status, recorded_by, recorded_at, scanned_at, notes, is_manual
		FROM attendances
		WHERE class_id = $1 AND attendance_date = $2
	`
	rows, err := r.db.Query(query, classID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		var scannedAt sql.NullTime
		var notes sql.NullString

		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.ClassID, &a.AttendanceDate, &a.Status,
			&a.RecordedBy, &a.RecordedAt, &scannedAt, &notes, &a.IsManual,
		); err != nil {
			return nil, err
		}

		if scannedAt.Valid {
			a.ScannedAt = &scannedAt.Time
		}
		if notes.Valid {
			str := notes.String
			a.Notes = &str
		}

		attendances = append(attendances, &a)
	}

	return attendances, nil
}
