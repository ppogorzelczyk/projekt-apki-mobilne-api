package users

import (
	"buymeagiftapi/internal/domain"
	"database/sql"
	"log/slog"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, username, phone, first_name, last_name, password_hash, created_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error(err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *repository) Create(user domain.User) (*domain.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	query := `SELECT id FROM users WHERE email = $1`
	var existingUserId int
	err = tx.QueryRow(query, user.Email).Scan(&existingUserId)

	if err == nil {
		tx.Rollback()
		return nil, nil
	}

	if err != sql.ErrNoRows {
		tx.Rollback()
		slog.Error(err.Error())
		return nil, err
	}

	user.CreatedAt = time.Now()
	query = `INSERT INTO users (email, username, password_hash, phone, first_name, last_name) 
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	err = tx.QueryRow(query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Phone,
		user.FirstName,
		user.LastName,
	).Scan(&id)

	if err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	user.Id = id

	err = tx.Commit()

	if err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	return &user, nil
}

func (r *repository) DeleteUser(userId int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, userId)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
