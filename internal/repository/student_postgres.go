package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

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
		INSERT INTO students (nisn, full_name, class_id, date_of_birth, gender, qr_code_data)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, is_active, created_at, updated_at
	`
	err := r.db.QueryRow(query, student.NISN, student.FullName, student.ClassID, student.DOB, student.Gender, student.QRCodeData).
		Scan(&student.ID, &student.IsActive, &student.CreatedAt, &student.UpdatedAt)
	if err != nil {
		slog.Error("Failed to insert student record", slog.Any("error", err), slog.String("nisn", student.NISN))
	}
	return err
}

func (r *studentPostgres) CreateBulk(students []*domain.Student) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO students (nisn, full_name, class_id, date_of_birth, gender, qr_code_data)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, is_active, created_at, updated_at
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, student := range students {
		err := stmt.QueryRow(student.NISN, student.FullName, student.ClassID, student.DOB, student.Gender, student.QRCodeData).
			Scan(&student.ID, &student.IsActive, &student.CreatedAt, &student.UpdatedAt)
		if err != nil {
			slog.Error("Failed to insert student record in bulk", slog.Any("error", err), slog.String("nisn", student.NISN))
			return err
		}
	}

	return tx.Commit()
}

func (r *studentPostgres) GetByID(id string) (*domain.Student, error) {
	var s domain.Student
	var dob, gender, photoURL sql.NullString

	query := `
		SELECT s.id, s.nisn, s.full_name, e.class_id, c.room_name, c.grade, c.section,
		       s.date_of_birth, s.gender, s.photo_url, s.qr_code_data,
		       s.is_active, s.created_at, s.updated_at
		FROM students s
		LEFT JOIN student_enrollments e ON s.id = e.student_id AND e.ended_at IS NULL
		LEFT JOIN classes c ON e.class_id = c.id
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`
	var roomName sql.NullString
	var grade sql.NullInt32
	var section sql.NullInt32
	var classID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.NISN, &s.FullName, &classID, &roomName, &grade, &section, &dob, &gender, &photoURL, &s.QRCodeData, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if classID.Valid {
		s.ClassID = classID.String
		if section.Valid {
			s.ClassName = fmt.Sprintf("Kelas %d %s %d", grade.Int32, roomName.String, section.Int32)
		} else if roomName.Valid {
			s.ClassName = fmt.Sprintf("Kelas %d %s", grade.Int32, roomName.String)
		}
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
		SELECT s.id, s.nisn, s.full_name, e.class_id, s.qr_code_data, s.is_active
		FROM students s
		LEFT JOIN student_enrollments e ON s.id = e.student_id AND e.ended_at IS NULL
		WHERE s.nisn = $1 AND s.is_active = true AND s.deleted_at IS NULL
	`
	var classID sql.NullString
	err := r.db.QueryRow(query, nisn).Scan(&s.ID, &s.NISN, &s.FullName, &classID, &s.QRCodeData, &s.IsActive)
	if classID.Valid {
		s.ClassID = classID.String
	}
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error("Failed to get student by NISN", slog.Any("error", err), slog.String("nisn", nisn))
		}
		return nil, err
	}
	return &s, nil
}

func (r *studentPostgres) GetByClassID(classID string) ([]*domain.Student, error) {
	query := `
		SELECT s.id, s.nisn, s.full_name, e.class_id, s.qr_code_data, s.is_active, s.created_at
		FROM students s
		JOIN student_enrollments e ON s.id = e.student_id AND e.ended_at IS NULL
		WHERE e.class_id = $1 AND s.is_active = true AND s.deleted_at IS NULL
		ORDER BY s.full_name ASC
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

func (r *studentPostgres) GetAll(teacherID string, isActive *bool, classID, search string, page, limit int) ([]*domain.Student, int, error) {
	offset := (page - 1) * limit
	whereClause := "WHERE s.deleted_at IS NULL"
	args := []interface{}{}
	argId := 1

	if isActive != nil {
		whereClause += fmt.Sprintf(" AND s.is_active = $%d", argId)
		args = append(args, *isActive)
		argId++
	}

	if classID != "" {
		whereClause += fmt.Sprintf(" AND e.class_id = $%d", argId)
		args = append(args, classID)
		argId++
	}

	if search != "" {
		whereClause += fmt.Sprintf(" AND (s.full_name ILIKE $%d OR s.nisn ILIKE $%d)", argId, argId)
		args = append(args, "%"+search+"%")
		argId++
	}

	if teacherID != "" {
		whereClause += fmt.Sprintf(" AND c.teacher_id = $%d", argId)
		args = append(args, teacherID)
		argId++
	}

	countQuery := "SELECT COUNT(*) FROM students s " +
		"LEFT JOIN student_enrollments e ON s.id = e.student_id AND e.ended_at IS NULL "
	if teacherID != "" {
		countQuery += "LEFT JOIN classes c ON e.class_id = c.id "
	}
	countQuery += whereClause
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT s.id, s.nisn, s.full_name, e.class_id, c.room_name, c.grade, c.section, s.date_of_birth, s.gender, s.photo_url, s.qr_code_data, s.is_active, s.created_at, s.updated_at " +
		"FROM students s LEFT JOIN student_enrollments e ON s.id = e.student_id AND e.ended_at IS NULL " +
		"LEFT JOIN classes c ON e.class_id = c.id " + whereClause + " " +
		"ORDER BY s.full_name ASC "

	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		var s domain.Student
		var classId sql.NullString
		var roomName sql.NullString
		var grade sql.NullInt32
		var section sql.NullInt32
		var dob, gender, photoURL sql.NullString
		if err := rows.Scan(&s.ID, &s.NISN, &s.FullName, &classId, &roomName, &grade, &section, &dob, &gender, &photoURL, &s.QRCodeData, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if classId.Valid {
			s.ClassID = classId.String
			if section.Valid {
				s.ClassName = fmt.Sprintf("Kelas %d %s %d", grade.Int32, roomName.String, section.Int32)
			} else if roomName.Valid {
				s.ClassName = fmt.Sprintf("Kelas %d %s", grade.Int32, roomName.String)
			}
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

func (r *studentPostgres) Update(student *domain.Student) error {
	query := `
		UPDATE students 
		SET full_name = $1, class_id = $2, date_of_birth = $3, gender = $4, is_active = $5, qr_code_data = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(query, student.FullName, student.ClassID, student.DOB, student.Gender, student.IsActive, student.QRCodeData, student.ID)
	return err
}

func (r *studentPostgres) SoftDelete(id string) error {
	query := `UPDATE students SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}
