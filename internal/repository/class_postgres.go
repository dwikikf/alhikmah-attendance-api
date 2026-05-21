package repository

import (
	"database/sql"
	"fmt"

	"alhikmah-attendance-api/internal/domain"
)

type classPostgres struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) domain.ClassRepository {
	return &classPostgres{db: db}
}

func (r *classPostgres) GetAll(teacherID string, academicYear string, offset, limit int) ([]*domain.Class, int, error) {
	query := `
		SELECT c.id, c.class_name, c.teacher_id, u.full_name, c.academic_year, c.capacity, c.description, c.created_at, c.updated_at,
		       (SELECT count(*) FROM students s WHERE s.class_id = c.id AND s.is_active = true) as student_count
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE 1=1
	`
	countQuery := `
		SELECT count(*) 
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE 1=1
	`

	args := []interface{}{}
	argId := 1

	if teacherID != "" {
		query += fmt.Sprintf(" AND c.teacher_id = $%d", argId)
		countQuery += fmt.Sprintf(" AND c.teacher_id = $%d", argId)
		args = append(args, teacherID)
		argId++
	}

	if academicYear != "" {
		query += fmt.Sprintf(" AND c.academic_year = $%d", argId)
		countQuery += fmt.Sprintf(" AND c.academic_year = $%d", argId)
		args = append(args, academicYear)
		argId++
	}

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY c.created_at DESC LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var classes []*domain.Class
	for rows.Next() {
		var c domain.Class
		var desc sql.NullString
		err := rows.Scan(
			&c.ID, &c.ClassName, &c.TeacherID, &c.TeacherName, &c.AcademicYear, &c.Capacity, &desc, &c.CreatedAt, &c.UpdatedAt, &c.StudentCount,
		)
		if err != nil {
			return nil, 0, err
		}
		if desc.Valid {
			c.Description = desc.String
		}
		classes = append(classes, &c)
	}

	return classes, total, nil
}

func (r *classPostgres) GetByID(id string) (*domain.Class, error) {
	var c domain.Class
	var desc sql.NullString
	query := `
		SELECT c.id, c.class_name, c.teacher_id, u.full_name, c.academic_year, c.capacity, c.description, c.created_at, c.updated_at,
		       (SELECT count(*) FROM students s WHERE s.class_id = c.id AND s.is_active = true) as student_count
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE c.id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.ClassName, &c.TeacherID, &c.TeacherName, &c.AcademicYear, &c.Capacity, &desc, &c.CreatedAt, &c.UpdatedAt, &c.StudentCount,
	)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		c.Description = desc.String
	}
	return &c, nil
}
