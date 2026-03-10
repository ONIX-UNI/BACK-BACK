package dto

import (
	"time"
)

// Tipos de trámite con SLA (Tiempos de respuesta)
type ProcedureType struct {
	ID              int16     `db:"id" json:"id"`
	Code            string    `db:"code" json:"code"`
	Name            string    `db:"name" json:"name"`
	SlaBusinessDays int       `db:"sla_business_days" json:"sla_business_days"`
	Alert48h        bool      `db:"alert_48h" json:"alert_48h"`
	Alert24h        bool      `db:"alert_24h" json:"alert_24h"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProcedureTypeRequest struct {
	Code            string `json:"code" db:"code"`
	Name            string `json:"name" db:"name"`
	SlaBusinessDays int    `json:"sla_business_days" db:"sla_business_day"`
	Alert48h        bool   `json:"alert_48h" db:"alert_48h"`
	Alert24h        bool   `json:"alert_24h" db:"alert_24h"`
}

//Tipos de documento aceptados en el sistema
type DocumentKind struct {
	ID        int16     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}