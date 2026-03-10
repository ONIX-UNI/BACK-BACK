package dto

type ServiceModality struct {
	ID       int16  `db:"id" json:"id"`
	Code     string `db:"code" json:"code"`
	Name     string `db:"name" json:"name"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

type ServiceModalityCreateRequest struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type ServiceModalityUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
