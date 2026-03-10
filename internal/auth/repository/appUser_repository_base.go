package repository

import (
	"database/sql"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppUserInstance struct {
	db *pgxpool.Pool
}

func NewAppUserInstance(db *pgxpool.Pool) *AppUserInstance {
	return &AppUserInstance{db: db}
}

type appUserScanner interface {
	Scan(dest ...any) error
}

func scanAppUser(scanner appUserScanner, user *dto.AppUser) error {
	var lastAccessAt sql.NullTime
	if err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.IsActive,
		&lastAccessAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return err
	}

	if lastAccessAt.Valid {
		lastAccessUTC := lastAccessAt.Time.UTC()
		user.LastAccessAt = &lastAccessUTC
	} else {
		user.LastAccessAt = nil
	}

	return nil
}
