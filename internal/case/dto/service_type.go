package dto

import "time"

type ServiceType struct {
	ID        int16     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ServiceTypeCreateRequest struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type ServiceTypeUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
