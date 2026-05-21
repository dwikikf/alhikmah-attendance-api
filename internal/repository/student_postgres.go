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
		SELECT id, nisn, full_name, class_id, date_of_birth, gender, photo_url, qr_code_data, is_active, created_at, updated_at
		FROM students WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.NISN, &s.FullName, &s.ClassID, &dob, &gender, &photoURL, &s.QRCodeData, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
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
		FROM students WHERE nisn = $1 AND is_active = true
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
		FROM students WHERE class_id = $1 AND is_active = true
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
