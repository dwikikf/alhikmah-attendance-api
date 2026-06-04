package repository

import (
	"database/sql"
	"fmt"


	"alhikmah-attendance-api/internal/domain"
)

type enrollmentPostgres struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) domain.EnrollmentRepository {
	return &enrollmentPostgres{db: db}
}

func (r *enrollmentPostgres) Enroll(e *domain.StudentEnrollment) error {
	query := `
		INSERT INTO student_enrollments (student_id, class_id, academic_year, enrolled_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, enrolled_at
	`
	err := r.db.QueryRow(query, e.StudentID, e.ClassID, e.AcademicYear).Scan(&e.ID, &e.EnrolledAt)
	if err != nil {
		return fmt.Errorf("failed to enroll student: %v", err)
	}
	return nil
}

func (r *enrollmentPostgres) GetActiveByStudentID(studentID string) (*domain.StudentEnrollment, error) {
	query := `
		SELECT e.id, e.student_id, s.full_name, e.class_id, e.academic_year, e.enrolled_at, e.ended_at, e.end_reason,
		       c.grade, c.room_name, c.section
		FROM student_enrollments e
		JOIN students s ON e.student_id = s.id
		JOIN classes c ON e.class_id = c.id
		WHERE e.student_id = $1 AND e.ended_at IS NULL
	`
	var e domain.StudentEnrollment
	var grade int
	var roomName string
	var section sql.NullInt32
	var endedAt sql.NullTime
	var endReason sql.NullString

	err := r.db.QueryRow(query, studentID).Scan(
		&e.ID, &e.StudentID, &e.StudentName, &e.ClassID, &e.AcademicYear, &e.EnrolledAt, &endedAt, &endReason,
		&grade, &roomName, &section,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no active enrollment
		}
		return nil, err
	}
	
	if endedAt.Valid {
		e.EndedAt = &endedAt.Time
	}
	if endReason.Valid {
		e.EndReason = &endReason.String
	}
	
	// Format class display name
	if section.Valid {
		e.ClassDisplay = fmt.Sprintf("Kelas %d %s %d", grade, roomName, section.Int32)
	} else {
		e.ClassDisplay = fmt.Sprintf("Kelas %d %s", grade, roomName)
	}

	return &e, nil
}

func (r *enrollmentPostgres) GetActiveByClassID(classID string) ([]*domain.StudentEnrollment, error) {
	query := `
		SELECT e.id, e.student_id, s.full_name, e.class_id, e.academic_year, e.enrolled_at,
		       c.grade, c.room_name, c.section
		FROM student_enrollments e
		JOIN students s ON e.student_id = s.id
		JOIN classes c ON e.class_id = c.id
		WHERE e.class_id = $1 AND e.ended_at IS NULL
		ORDER BY s.full_name ASC
	`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*domain.StudentEnrollment
	for rows.Next() {
		var e domain.StudentEnrollment
		var grade int
		var roomName string
		var section sql.NullInt32

		err := rows.Scan(
			&e.ID, &e.StudentID, &e.StudentName, &e.ClassID, &e.AcademicYear, &e.EnrolledAt,
			&grade, &roomName, &section,
		)
		if err != nil {
			return nil, err
		}

		if section.Valid {
			e.ClassDisplay = fmt.Sprintf("Kelas %d %s %d", grade, roomName, section.Int32)
		} else {
			e.ClassDisplay = fmt.Sprintf("Kelas %d %s", grade, roomName)
		}
		
		enrollments = append(enrollments, &e)
	}

	return enrollments, nil
}

func (r *enrollmentPostgres) GetHistoryByStudentID(studentID string) ([]*domain.StudentEnrollment, error) {
	query := `
		SELECT e.id, e.student_id, s.full_name, e.class_id, e.academic_year, e.enrolled_at, e.ended_at, e.end_reason,
		       c.grade, c.room_name, c.section
		FROM student_enrollments e
		JOIN students s ON e.student_id = s.id
		JOIN classes c ON e.class_id = c.id
		WHERE e.student_id = $1
		ORDER BY e.enrolled_at DESC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*domain.StudentEnrollment
	for rows.Next() {
		var e domain.StudentEnrollment
		var grade int
		var roomName string
		var section sql.NullInt32
		var endedAt sql.NullTime
		var endReason sql.NullString

		err := rows.Scan(
			&e.ID, &e.StudentID, &e.StudentName, &e.ClassID, &e.AcademicYear, &e.EnrolledAt, &endedAt, &endReason,
			&grade, &roomName, &section,
		)
		if err != nil {
			return nil, err
		}

		if endedAt.Valid {
			e.EndedAt = &endedAt.Time
		}
		if endReason.Valid {
			e.EndReason = &endReason.String
		}

		if section.Valid {
			e.ClassDisplay = fmt.Sprintf("Kelas %d %s %d", grade, roomName, section.Int32)
		} else {
			e.ClassDisplay = fmt.Sprintf("Kelas %d %s", grade, roomName)
		}
		
		enrollments = append(enrollments, &e)
	}

	return enrollments, nil
}

func (r *enrollmentPostgres) EndEnrollment(studentID, classID string, reason string) error {
	query := `
		UPDATE student_enrollments
		SET ended_at = NOW(), end_reason = $3
		WHERE student_id = $1 AND class_id = $2 AND ended_at IS NULL
	`
	_, err := r.db.Exec(query, studentID, classID, reason)
	return err
}

func (r *enrollmentPostgres) BulkEnroll(items []domain.PromoteItem, academicYear string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	query := `
		INSERT INTO student_enrollments (student_id, class_id, academic_year, enrolled_at)
		VALUES ($1, $2, $3, NOW())
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count := 0
	for _, item := range items {
		// Close existing enrollment
		closeQuery := `
			UPDATE student_enrollments
			SET ended_at = NOW(), end_reason = 'promoted'
			WHERE student_id = $1 AND ended_at IS NULL
		`
		_, err = tx.Exec(closeQuery, item.StudentID)
		if err != nil {
			return 0, err
		}

		// Insert new enrollment
		_, err = stmt.Exec(item.StudentID, item.TargetClassID, academicYear)
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}
