package dto

import (
	"time"
)

type CatalogItem struct {
	ID        int16     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type CreateCatalogItemRequest struct {
	Code     string `json:"code" db:"code"`
	Name     string `json:"name" db:"name"`
	IsActive *bool  `json:"is_active,omitempty" db:"is_active"`
}

type UpdateCatalogItemRequest struct {
	Name     *string `json:"name,omitempty" db:"name"`
	IsActive *bool   `json:"is_active,omitempty" db:"is_activate"`
}

type DerivationEntity struct {
	CatalogItem
	ContactInfo *string `db:"contact_info" json:"contact_info,omitempty"`
}

type CreateDerivationEntityRequest struct {
	CreateCatalogItemRequest
	ContactInfo *string `json:"contact_info,omitempty" db:"contact_info"`
}