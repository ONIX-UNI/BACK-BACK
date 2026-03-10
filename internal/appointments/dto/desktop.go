package dto

import (
	"time"

	"github.com/google/uuid"
)

type Desktop struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Code     string    `db:"code" json:"code"`
	Name     string    `db:"name" json:"name"`
	Location *string   `db:"location" json:"location,omitempty"`
	IsActive bool      `db:"is_active" json:"is_active"`

	CurrentUserID *uuid.UUID `db:"current_user_id" json:"current_user_id,omitempty"`
	AssignedAt    *time.Time `db:"assigned_at" json:"assigned_at,omitempty"`
	AssignedBy    *uuid.UUID `db:"assigned_by" json:"assigned_by,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateDesktopRequest struct {
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Location *string `json:"location,omitempty"`
}

type CreateDesktopResponse struct {
	ID string `json:"id"`
}

// ✅ NUEVO: PATCH /escritorios/:code (flexible)
type UpdateDesktopRequest struct {
	Name     *string    `json:"name,omitempty"`
	IsActive *bool      `json:"is_active,omitempty"`
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Reason   *string    `json:"reason,omitempty"`
}

type UpdateDesktopResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}