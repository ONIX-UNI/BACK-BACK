package auth

import (
	"context"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/handlers"
	"github.com/DuvanRozoParra/sicou/internal/auth/repository"
	"github.com/DuvanRozoParra/sicou/internal/auth/service"
	"github.com/DuvanRozoParra/sicou/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type ModuleConfig struct {
	JWTSecret             string
	TokenTTL              time.Duration
	PasswordResetTokenTTL time.Duration
	PasswordResetURL      string
	CookieName            string
	CookiePath            string
	CookieDomain          string
	CookieSecure          bool
	CookieSameSite        string
	SuperAdminEmail       string
	SuperAdminDisplayName string
	SuperAdminPassword    string
}

type Module struct {
	HandlerAppUser *handlers.AppUserHandler
}

func NewModule(ctx context.Context, cfg ModuleConfig) (*Module, error) {
	repo_appUser := repository.NewAppUserInstance(database.DB)

	svc_appUser := service.NewAppUserService(
		repo_appUser,
		cfg.JWTSecret,
		cfg.TokenTTL,
		service.AppUserServiceConfig{
			PasswordResetTokenTTL: cfg.PasswordResetTokenTTL,
			PasswordResetURL:      cfg.PasswordResetURL,
		},
	)

	if err := svc_appUser.EnsureAuthStorage(ctx); err != nil {
		return nil, err
	}

	if err := svc_appUser.EnsureSingleSuperAdmin(ctx, dto.BootstrapSuperAdminRequest{
		Email:       cfg.SuperAdminEmail,
		DisplayName: cfg.SuperAdminDisplayName,
		Password:    cfg.SuperAdminPassword,
	}); err != nil {
		return nil, err
	}

	handler_appUser := handlers.NewAppUserHandler(svc_appUser, handlers.SessionCookieConfig{
		Name:     cfg.CookieName,
		Path:     cfg.CookiePath,
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		SameSite: cfg.CookieSameSite,
	})

	return &Module{
		HandlerAppUser: handler_appUser,
	}, nil
}

func (m *Module) RequireSession() fiber.Handler {
	return m.HandlerAppUser.RequireSession
}

func (m *Module) OptionalSession() fiber.Handler {
	return m.HandlerAppUser.OptionalSession
}
