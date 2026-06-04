package repository

import (
	"database/sql"
	"log/slog"
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
		INSERT INTO attendances (student_id, class_id, attendance_date, subject, status, recorded_by, scanned_at, is_manual, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, recorded_at
	`
	err := r.db.QueryRow(
		query,
		attendance.StudentID,
		attendance.ClassID,
		attendance.AttendanceDate,
		attendance.Subject,
		attendance.Status,
		attendance.RecordedBy,
		attendance.ScannedAt,
		attendance.IsManual,
		attendance.Notes,
	).Scan(&attendance.ID, &attendance.RecordedAt)

	if err != nil {
		slog.Error("Failed to insert attendance record", slog.Any("error", err), slog.String("student_id", attendance.StudentID))
	}
	return err
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
	`
	_, err = tx.Exec(updateQuery, attendance.Status, attendance.RecordedBy, attendance.ID)
	if err != nil {
		slog.Error("Failed to update attendance record", slog.Any("error", err), slog.String("attendance_id", attendance.ID))
		return err
	}

	// Insert audit
	auditQuery := `
		INSERT INTO attendance_audits (attendance_id, old_status, new_status, changed_by, reason)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(auditQuery, audit.AttendanceID, audit.OldStatus, audit.NewStatus, audit.ChangedBy, audit.Reason)
	if err != nil {
		slog.Error("Failed to insert attendance audit record", slog.Any("error", err), slog.String("attendance_id", audit.AttendanceID))
		return err
	}

	return tx.Commit()
}

func (r *attendancePostgres) GetByClassAndDate(classID string, date time.Time, subject string) ([]*domain.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, attendance_date, subject, status, recorded_by, recorded_at, scanned_at, notes, is_manual
		FROM attendances
		WHERE class_id = $1 AND attendance_date = $2
	`
	args := []interface{}{classID, date}
	if subject == "" {
		query += " AND subject IS NULL"
	} else {
		query += " AND subject = $3"
		args = append(args, subject)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		slog.Error("Failed to get attendances by class and date", slog.Any("error", err), slog.String("class_id", classID))
		return nil, err
	}
	defer rows.Close()

	var attendances []*domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		var scannedAt sql.NullTime
		var notes sql.NullString

		var subjectNull sql.NullString

		if err := rows.Scan(
			&a.ID, &a.StudentID, &a.ClassID, &a.AttendanceDate, &subjectNull, &a.Status,
			&a.RecordedBy, &a.RecordedAt, &scannedAt, &notes, &a.IsManual,
		); err != nil {
			return nil, err
		}

		if subjectNull.Valid {
			str := subjectNull.String
			a.Subject = &str
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

func (r *attendancePostgres) GetByID(id string) (*domain.Attendance, error) {
	query := `
		SELECT id, student_id, class_id, attendance_date, subject, status, recorded_by, recorded_at, scanned_at, notes, is_manual
		FROM attendances
		WHERE id = $1
	`
	var a domain.Attendance
	var scannedAt sql.NullTime
	var notes sql.NullString
	var subjectNull sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.StudentID, &a.ClassID, &a.AttendanceDate, &subjectNull, &a.Status,
		&a.RecordedBy, &a.RecordedAt, &scannedAt, &notes, &a.IsManual,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		slog.Error("Failed to get attendance by ID", slog.Any("error", err), slog.String("attendance_id", id))
		return nil, err
	}

	if scannedAt.Valid {
		a.ScannedAt = &scannedAt.Time
	}
	if notes.Valid {
		str := notes.String
		a.Notes = &str
	}
	if subjectNull.Valid {
		str := subjectNull.String
		a.Subject = &str
	}

	return &a, nil
}
