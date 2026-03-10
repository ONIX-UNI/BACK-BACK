package models

import (
	"context"
	"errors"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
)

var ErrRoleNotFound = errors.New("role not found")

type AppUserRepository interface {
	Create(ctx context.Context, req dto.CreateAppUserRequest) (*dto.AppUser, error)

	GetByID(ctx context.Context, id uuid.UUID) (*dto.AppUser, error)
	GetByEmail(ctx context.Context, email string) (*dto.AppUser, error)
	EnqueuePasswordResetEmail(ctx context.Context, toEmail string, subject string, body string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
	RoleExists(ctx context.Context, roleCode string) (bool, error)
	GetByDisplayName(ctx context.Context, displayName string, limit, offset int) (*dto.ListAppUsersResponse, error)
	ListOptions(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error)
	ListWithFilters(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error)

	List(ctx context.Context, limit, offset int) (*dto.ListAppUsersResponse, error)

	Update(ctx context.Context, id uuid.UUID, req dto.UpdateAppUserRequest) (*dto.AppUser, error)
	ReplacePrimaryRole(ctx context.Context, userID uuid.UUID, roleCode string) error
	UpdateLastAccess(ctx context.Context, id uuid.UUID, accessedAt time.Time) error

	EnsureSingleSuperAdmin(ctx context.Context, email, displayName, passwordHash string) error
	EnsureAuthStorage(ctx context.Context) error
	IsTokenRevoked(ctx context.Context, tokenSignature string) (bool, error)
	RevokeToken(ctx context.Context, tokenSignature string, expiresAt time.Time) error
	CreateUserAuditLog(ctx context.Context, entry dto.UserAuditLogEntry) error
	ListUserAuditLog(ctx context.Context, req dto.ListUserAuditLogRequest) (*dto.ListUserAuditLogResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
