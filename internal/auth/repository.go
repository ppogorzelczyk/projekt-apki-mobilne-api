package auth

import (
	"database/sql"
	"log/slog"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) BlacklistToken(token string) error {
	query := `INSERT INTO token_blacklists (token) VALUES ($1)`
	_, err := r.db.Exec(query, token)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) IsTokenBlacklisted(token string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM token_blacklists WHERE token = $1)`
	var isBlacklisted bool
	err := r.db.QueryRow(query, token).Scan(&isBlacklisted)

	if err != nil {
		slog.Error(err.Error())
		return false, err
	}

	return isBlacklisted, nil
}
