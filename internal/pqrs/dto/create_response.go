package dto

import "time"

type CreatePQRSResponse struct {
	ID        string    `json:"id"`
	Radicado  string    `json:"radicado"`
	Estado    string    `json:"estado"`
	CreatedAt time.Time `json:"createdAt"`
}
