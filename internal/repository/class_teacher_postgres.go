package repository

import (
	"database/sql"
	"fmt"

	"alhikmah-attendance-api/internal/domain"
)

type classTeacherPostgres struct {
	db *sql.DB
}

func NewClassTeacherRepository(db *sql.DB) domain.ClassTeacherRepository {
	return &classTeacherPostgres{db: db}
}

func (r *classTeacherPostgres) Assign(ct *domain.ClassTeacher) error {
	query := `
		INSERT INTO class_teachers (teacher_id, class_id, academic_year, subject, role, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, ct.TeacherID, ct.ClassID, ct.AcademicYear, ct.Subject, ct.Role).Scan(&ct.ID, &ct.CreatedAt)
}

func (r *classTeacherPostgres) Unassign(teacherID, classID, subject string) error {
	query := `
		DELETE FROM class_teachers
		WHERE teacher_id = $1 AND class_id = $2 AND subject = $3
	`
	_, err := r.db.Exec(query, teacherID, classID, subject)
	return err
}

func (r *classTeacherPostgres) GetByClassID(classID string) ([]*domain.ClassTeacher, error) {
	// Include homeroom teacher from `classes`
	query := `
		SELECT u.id, u.full_name, c.id, c.academic_year, 'homeroom' as role, '' as subject, c.room_name, c.grade, c.section, c.created_at
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE c.id = $1
		UNION ALL
		SELECT ct.teacher_id, u.full_name, ct.class_id, ct.academic_year, ct.role, COALESCE(ct.subject, ''), c.room_name, c.grade, c.section, ct.created_at
		FROM class_teachers ct
		JOIN users u ON ct.teacher_id = u.id
		JOIN classes c ON ct.class_id = c.id
		WHERE ct.class_id = $1
	`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ClassTeacher
	for rows.Next() {
		var ct domain.ClassTeacher
		var roomName string
		var grade int
		var section sql.NullInt32

		if err := rows.Scan(
			&ct.TeacherID, &ct.TeacherName, &ct.ClassID, &ct.AcademicYear, &ct.Role, &ct.Subject,
			&roomName, &grade, &section, &ct.CreatedAt,
		); err != nil {
			return nil, err
		}

		if section.Valid {
			ct.ClassDisplay = fmt.Sprintf("Kelas %d %s %d", grade, roomName, section.Int32)
		} else {
			ct.ClassDisplay = fmt.Sprintf("Kelas %d %s", grade, roomName)
		}

		list = append(list, &ct)
	}

	return list, nil
}

func (r *classTeacherPostgres) GetByTeacherID(teacherID string) ([]*domain.ClassTeacher, error) {
	query := `
		SELECT ct.id, ct.teacher_id, u.full_name, ct.class_id, ct.academic_year, ct.role, COALESCE(ct.subject, ''), c.room_name, c.grade, c.section, ct.created_at
		FROM class_teachers ct
		JOIN users u ON ct.teacher_id = u.id
		JOIN classes c ON ct.class_id = c.id
		WHERE ct.teacher_id = $1
	`
	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ClassTeacher
	for rows.Next() {
		var ct domain.ClassTeacher
		var roomName string
		var grade int
		var section sql.NullInt32

		if err := rows.Scan(
			&ct.ID, &ct.TeacherID, &ct.TeacherName, &ct.ClassID, &ct.AcademicYear, &ct.Role, &ct.Subject,
			&roomName, &grade, &section, &ct.CreatedAt,
		); err != nil {
			return nil, err
		}

		if section.Valid {
			ct.ClassDisplay = fmt.Sprintf("Kelas %d %s %d", grade, roomName, section.Int32)
		} else {
			ct.ClassDisplay = fmt.Sprintf("Kelas %d %s", grade, roomName)
		}

		list = append(list, &ct)
	}

	return list, nil
}

func (r *classTeacherPostgres) GetSubjectAssignments(teacherID string) ([]*domain.ClassTeacher, error) {
	query := `
		SELECT ct.id, ct.teacher_id, u.full_name, ct.class_id, ct.academic_year, ct.role, COALESCE(ct.subject, ''), c.room_name, c.grade, c.section, ct.created_at
		FROM class_teachers ct
		JOIN users u ON ct.teacher_id = u.id
		JOIN classes c ON ct.class_id = c.id
		WHERE ct.teacher_id = $1 AND ct.role = 'subject_teacher'
	`
	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ClassTeacher
	for rows.Next() {
		var ct domain.ClassTeacher
		var roomName string
		var grade int
		var section sql.NullInt32

		if err := rows.Scan(
			&ct.ID, &ct.TeacherID, &ct.TeacherName, &ct.ClassID, &ct.AcademicYear, &ct.Role, &ct.Subject,
			&roomName, &grade, &section, &ct.CreatedAt,
		); err != nil {
			return nil, err
		}

		if section.Valid {
			ct.ClassDisplay = fmt.Sprintf("Kelas %d %s %d", grade, roomName, section.Int32)
		} else {
			ct.ClassDisplay = fmt.Sprintf("Kelas %d %s", grade, roomName)
		}

		list = append(list, &ct)
	}

	return list, nil
}
