package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/auth/service"
)

type SessionCookieConfig struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	SameSite string
}

type AppUserHandler struct {
	service           *service.AppUserService
	sessionCookieConf SessionCookieConfig
}

const maxBulkUploadBytes int64 = 10 * 1024 * 1024

func NewAppUserHandler(service *service.AppUserService, sessionCookieConf SessionCookieConfig) *AppUserHandler {
	return &AppUserHandler{
		service:           service,
		sessionCookieConf: normalizeSessionCookieConfig(sessionCookieConf),
	}
}
