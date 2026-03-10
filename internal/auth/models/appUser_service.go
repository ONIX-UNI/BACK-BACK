package models

import (
	"context"
	"io"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
)

type IAppUserService interface {
	Create(ctx context.Context, req dto.CreateAppUserRequest, actorUserID *uuid.UUID) (*dto.AppUser, error)
	BulkCreate(ctx context.Context, file io.Reader, actorUserID *uuid.UUID) (*dto.BulkCreateAppUserResponse, error)

	GetByID(ctx context.Context, id uuid.UUID) (*dto.AppUser, error)

	GetByEmail(ctx context.Context, email string) (*dto.AppUser, error)

	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error

	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error)
	Me(ctx context.Context, rawToken string) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, rawToken string) (*dto.LoginResponse, string, error)
	Logout(ctx context.Context, rawToken string) error

	GetByDisplayName(ctx context.Context, displayName string, limit, offset int) (*dto.ListAppUsersResponse, error)
	ListOptions(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error)
	ListWithFilters(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error)

	List(ctx context.Context, limit, offset int) (*dto.ListAppUsersResponse, error)
	ListAuditLog(ctx context.Context, req dto.ListUserAuditLogRequest) (*dto.ListUserAuditLogResponse, error)

	Update(ctx context.Context, id uuid.UUID, req dto.UpdateAppUserRequest, actorUserID *uuid.UUID) (*dto.AppUser, error)

	EnsureAuthStorage(ctx context.Context) error
	EnsureSingleSuperAdmin(ctx context.Context, req dto.BootstrapSuperAdminRequest) error

	Delete(ctx context.Context, id uuid.UUID, actorUserID *uuid.UUID) error
}
