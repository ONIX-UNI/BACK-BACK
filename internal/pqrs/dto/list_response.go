package dto

import "time"

type PQRSListItem struct {
	ID           string    `json:"id"`
	Radicado     string    `json:"radicado"`
	Tipo         string    `json:"tipo"`
	Asunto       string    `json:"asunto"`
	Ciudadano    string    `json:"ciudadano"`
	Correo       string    `json:"correo"`
	Estado       string    `json:"estado"`
	Responsable  string    `json:"responsable"`
	FechaIngreso time.Time `json:"fechaIngreso"`
	FechaLimite  time.Time `json:"fechaLimite"`
}

type PQRSListMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type ListPQRSResponse struct {
	Data []PQRSListItem `json:"data"`
	Meta PQRSListMeta   `json:"meta"`
}
