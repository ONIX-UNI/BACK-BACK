package dto

import "time"

type CreateAsesoriaResponse struct {
	ID             string    `json:"id"`
	PreturnoNumber string    `json:"preturnoNumber"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}
