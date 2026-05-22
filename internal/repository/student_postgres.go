package repository

import (
	"database/sql"

	"alhikmah-attendance-api/internal/domain"
)

type studentPostgres struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) domain.StudentRepository {
	return &studentPostgres{db: db}
}

func (r *studentPostgres) Create(student *domain.Student) error {
	query := `
		INSERT INTO students (nisn, full_name, class_id, qr_code_data)
		VALUES ($1, $2, $3, $4)
		RETURNING id, is_active, created_at, updated_at
	`
	return r.db.QueryRow(query, student.NISN, student.FullName, student.ClassID, student.QRCodeData).
		Scan(&student.ID, &student.IsActive, &student.CreatedAt, &student.UpdatedAt)
}

func (r *studentPostgres) GetByID(id string) (*domain.Student, error) {
	var s domain.Student
	var dob, gender, photoURL sql.NullString

	query := `
		SELECT s.id, s.nisn, s.full_name, s.class_id, c.class_name,
		       s.date_of_birth, s.gender, s.photo_url, s.qr_code_data,
		       s.is_active, s.created_at, s.updated_at
		FROM students s
		JOIN classes c ON s.class_id = c.id
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.NISN, &s.FullName, &s.ClassID, &s.ClassName, &dob, &gender, &photoURL, &s.QRCodeData, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if dob.Valid {
		str := dob.String
		s.DOB = &str
	}
	if gender.Valid {
		str := gender.String
		s.Gender = &str
	}
	if photoURL.Valid {
		str := photoURL.String
		s.PhotoURL = &str
	}

	return &s, nil
}

func (r *studentPostgres) GetByNISN(nisn string) (*domain.Student, error) {
	var s domain.Student
	query := `
		SELECT id, nisn, full_name, class_id, qr_code_data, is_active
		FROM students WHERE nisn = $1 AND is_active = true AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, nisn).Scan(&s.ID, &s.NISN, &s.FullName, &s.ClassID, &s.QRCodeData, &s.IsActive)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentPostgres) GetByClassID(classID string) ([]*domain.Student, error) {
	query := `
		SELECT id, nisn, full_name, class_id, qr_code_data, is_active, created_at
		FROM students WHERE class_id = $1 AND is_active = true AND deleted_at IS NULL
		ORDER BY full_name ASC
	`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(&s.ID, &s.NISN, &s.FullName, &s.ClassID, &s.QRCodeData, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, err
		}
		students = append(students, &s)
	}

	return students, nil
}

func (r *studentPostgres) GetAll(isActive *bool, page, limit int) ([]*domain.Student, int, error) {
	offset := (page - 1) * limit
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argId := 1

	if isActive != nil {
		whereClause += " AND is_active = $1"
		args = append(args, *isActive)
		argId++
	}

	countQuery := "SELECT COUNT(*) FROM students " + whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, nisn, full_name, class_id, date_of_birth, gender, photo_url, qr_code_data, is_active, created_at, updated_at " +
		"FROM students " + whereClause + " " +
		"ORDER BY full_name ASC "

	if argId == 1 {
		query += "LIMIT $1 OFFSET $2"
	} else {
		query += "LIMIT $2 OFFSET $3"
	}
	
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		var s domain.Student
		var dob, gender, photoURL sql.NullString
		if err := rows.Scan(&s.ID, &s.NISN, &s.FullName, &s.ClassID, &dob, &gender, &photoURL, &s.QRCodeData, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if dob.Valid {
			str := dob.String
			s.DOB = &str
		}
		if gender.Valid {
			str := gender.String
			s.Gender = &str
		}
		if photoURL.Valid {
			str := photoURL.String
			s.PhotoURL = &str
		}
		students = append(students, &s)
	}

	return students, total, nil
}

func (r *studentPostgres) SoftDelete(id string) error {
	query := `UPDATE students SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}
