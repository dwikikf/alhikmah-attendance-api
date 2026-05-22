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
		       (SELECT count(*) FROM students s WHERE s.class_id = c.id AND s.is_active = true AND s.deleted_at IS NULL) as student_count
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE c.deleted_at IS NULL
	`
	countQuery := `
		SELECT count(*) 
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE c.deleted_at IS NULL
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
		       (SELECT count(*) FROM students s WHERE s.class_id = c.id AND s.is_active = true AND s.deleted_at IS NULL) as student_count
		FROM classes c
		JOIN users u ON c.teacher_id = u.id
		WHERE c.id = $1 AND c.deleted_at IS NULL
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

func (r *classPostgres) Create(class *domain.Class) error {
	query := `
		INSERT INTO classes (class_name, teacher_id, academic_year, capacity, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, class.ClassName, class.TeacherID, class.AcademicYear, class.Capacity, class.Description).Scan(
		&class.ID, &class.CreatedAt, &class.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *classPostgres) Update(class *domain.Class) error {
	query := `
		UPDATE classes 
		SET class_name = $1, teacher_id = $2, capacity = $3, description = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRow(query, class.ClassName, class.TeacherID, class.Capacity, class.Description, class.ID).Scan(
		&class.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("class not found")
		}
		return err
	}
	return nil
}

func (r *classPostgres) SoftDelete(id string) error {
	query := `UPDATE classes SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}
