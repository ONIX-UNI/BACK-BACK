package service

import (
	"errors"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/models"
)

type AppUserService struct {
	repo                     models.AppUserRepository
	jwtSecret                string
	tokenTTL                 time.Duration
	passwordResetTokenTTL    time.Duration
	passwordResetFrontendURL string
}

type AppUserServiceConfig struct {
	PasswordResetTokenTTL time.Duration
	PasswordResetURL      string
}

func NewAppUserService(
	repo models.AppUserRepository,
	jwtSecret string,
	tokenTTL time.Duration,
	configs ...AppUserServiceConfig,
) *AppUserService {
	if tokenTTL <= 0 {
		tokenTTL = 2 * time.Hour
	}

	cfg := AppUserServiceConfig{
		PasswordResetTokenTTL: 30 * time.Minute,
		PasswordResetURL:      "",
	}
	if len(configs) > 0 {
		if configs[0].PasswordResetTokenTTL > 0 {
			cfg.PasswordResetTokenTTL = configs[0].PasswordResetTokenTTL
		}
		cfg.PasswordResetURL = strings.TrimSpace(configs[0].PasswordResetURL)
	}

	return &AppUserService{
		repo:                     repo,
		jwtSecret:                strings.TrimSpace(jwtSecret),
		tokenTTL:                 tokenTTL,
		passwordResetTokenTTL:    cfg.PasswordResetTokenTTL,
		passwordResetFrontendURL: cfg.PasswordResetURL,
	}
}

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyUsed    = errors.New("email already in use")
	ErrRoleNotFound        = errors.New("role not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserInactive        = errors.New("user is inactive")
	ErrTokenNotConfigured  = errors.New("auth token secret is not configured")
	ErrPasswordHashFailure = errors.New("failed to hash password")
	ErrFileTooLarge        = errors.New("file is too large")
	ErrSessionInvalid      = errors.New("invalid session")
	ErrSessionExpired      = errors.New("session expired")
	ErrSessionRevoked      = errors.New("session revoked")
	ErrPasswordResetTokenInvalid = errors.New("invalid password reset token")
	ErrPasswordResetTokenExpired = errors.New("expired password reset token")
)
