package dto

import (
	"time"

	"github.com/google/uuid"
)

type AppUser struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	DisplayName  string     `json:"display_name" db:"display_name"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         string     `json:"role,omitempty" db:"-"`
	Roles        []string   `json:"roles,omitempty" db:"-"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastAccessAt *time.Time `json:"last_access_at,omitempty" db:"last_access_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateAppUserRequest struct {
	Email        string `json:"email" validate:"required,email"`
	DisplayName  string `json:"display_name" validate:"required"`
	PasswordHash string `json:"password_hash" validate:"required"`
	Role         string `json:"role" validate:"required"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

type CreateAppUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type BulkCreateSummary struct {
	Total   int `json:"total"`
	Created int `json:"created"`
	Failed  int `json:"failed"`
}

type BulkCreateResultItem struct {
	Line    int    `json:"line"`
	Email   string `json:"email"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type BulkCreateAppUserResponse struct {
	Summary BulkCreateSummary      `json:"summary"`
	Results []BulkCreateResultItem `json:"results"`
}

type UpdateAppUserRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	Role        *string `json:"role,omitempty"`
}

type UpdatePasswordRequest struct {
	PasswordHash string `json:"password_hash" validate:"required"`
}

type DeleteAppUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetAppUserByIDRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetAppUserByEmailRequest struct {
	Email string `json:"email"`
}

type GetAppUserByDisplayNameRequest struct {
	DisplayName string `json:"display_name"`
}

type ListAppUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListAppUsersFilterRequest struct {
	Query  string `json:"q"`
	Role   string `json:"role"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ListAppUsersResponse struct {
	Items       []AppUser `json:"items"`
	Total       int64     `json:"total"`
	ActiveCount int64     `json:"active_count"`
	Limit       int       `json:"limit"`
	Offset      int       `json:"offset"`
	HasNext     bool      `json:"has_next"`
}

type ListAppUserOptionsRequest struct {
	Roles    []string `json:"roles"`
	Search   string   `json:"search"`
	IsActive bool     `json:"is_active"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

type AppUserOptionItem struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
}

type ListAppUserOptionsResponse struct {
	Items   []AppUserOptionItem `json:"items"`
	Total   int64               `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	HasNext bool                `json:"has_next"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Roles       []string  `json:"roles"`
}

type LoginResponse struct {
	ExpiresAt time.Time         `json:"expires_at"`
	ExpiresIn int64             `json:"expires_in"`
	User      LoginUserResponse `json:"user"`
}

const (
	UserAuditActionCreated       = "USER_CREATED"
	UserAuditActionUpdated       = "USER_UPDATED"
	UserAuditActionDeleted       = "USER_DELETED"
	UserAuditActionActivated     = "USER_ACTIVATED"
	UserAuditActionDeactivated   = "USER_DEACTIVATED"
	UserAuditActionRoleChanged   = "USER_ROLE_CHANGED"
	UserAuditActionPasswordReset = "USER_PASSWORD_RESET"
)

type UserAuditActor struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserAuditTargetUser struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserAuditMetadata struct {
	Before map[string]any `json:"before,omitempty"`
	After  map[string]any `json:"after,omitempty"`
}

type UserAuditItem struct {
	ID          string              `json:"id"`
	CreatedAt   time.Time           `json:"createdAt"`
	Actor       UserAuditActor      `json:"actor"`
	Action      string              `json:"action"`
	ActionLabel string              `json:"actionLabel"`
	TargetUser  UserAuditTargetUser `json:"targetUser"`
	Detail      string              `json:"detail"`
	Metadata    UserAuditMetadata   `json:"metadata"`
}

type ListUserAuditLogRequest struct {
	Limit        int
	Offset       int
	From         *time.Time
	To           *time.Time
	Action       string
	ActorUserID  *uuid.UUID
	TargetUserID *uuid.UUID
	Search       string
}

type ListUserAuditLogResponse struct {
	Items  []UserAuditItem `json:"items"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

type UserAuditLogEntry struct {
	ActorUserID     *uuid.UUID
	ActorName       string
	ActorEmail      string
	Action          string
	TargetUserID    *uuid.UUID
	TargetUserName  string
	TargetUserEmail string
	Detail          string
	Metadata        UserAuditMetadata
}

type BootstrapSuperAdminRequest struct {
	Email       string
	DisplayName string
	Password    string
}
