package repository

import (
	"database/sql"
	"fmt"

	"alhikmah-attendance-api/internal/domain"
)

type userPostgres struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userPostgres{db: db}
}

func (r *userPostgres) GetByID(id string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, role, is_active, created_at, updated_at, last_login 
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName, &user.Role,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPostgres) GetAll(role string, isActive *bool, offset, limit int) ([]*domain.User, int, error) {
	query := `SELECT id, username, email, full_name, role, is_active, created_at, updated_at, last_login FROM users WHERE deleted_at IS NULL`
	countQuery := `SELECT count(*) FROM users WHERE deleted_at IS NULL`
	
	args := []interface{}{}
	argId := 1

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argId)
		countQuery += fmt.Sprintf(" AND role = $%d", argId)
		args = append(args, role)
		argId++
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argId)
		countQuery += fmt.Sprintf(" AND is_active = $%d", argId)
		args = append(args, *isActive)
		argId++
	}

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FullName, &u.Role,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.LastLogin,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}

	return users, total, nil
}

func (r *userPostgres) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at, updated_at
	`
	return r.db.QueryRow(
		query, user.Username, user.Email, user.PasswordHash, user.FullName, user.Role,
	).Scan(&user.ID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userPostgres) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET full_name = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`
	return r.db.QueryRow(query, user.FullName, user.Email, user.ID).Scan(&user.UpdatedAt)
}

func (r *userPostgres) SoftDelete(id string) error {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

