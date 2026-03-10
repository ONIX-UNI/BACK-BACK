package service

import (
	"context"
	"mime/multipart"
	"time"

	documentsdto "github.com/DuvanRozoParra/sicou/internal/documents/dto"
)

const (
	statusPublicRecibido = "RECIBIDO"
)

type FileObjectService interface {
	Create(ctx context.Context, file multipart.File, size int64, filename string, contentType string) (*documentsdto.FileObject, error)
	Delete(ctx context.Context, id string) error
}

type CreatePQRSInput struct {
	QueryType                 string
	PersonType                string
	DocumentType              string
	DocumentOrTaxID           string
	FirstName                 string
	MiddleName                string
	FirstLastName             string
	SecondLastName            string
	Gender                    string
	Address                   string
	NeighborhoodArea          string
	AllowsElectronicResponse  string
	Email                     string
	Phone                     string
	PopulationGroup           string
	OtherPopulationGroup      string
	RequestDescription        string
	ResponseChannel           string
	RequestType               string
	RequestAgainstStudent     string
	ResponsibleStudentName    string
	ResponsibleStudentProgram string
	StudentCaseDescription    string
	AcceptsDataProcessing     string
	SubmittedAt               string
	Attachments               []*multipart.FileHeader
}

type CreatePQRSResult struct {
	ID        string
	Radicado  string
	Estado    string
	CreatedAt time.Time
}

type ListPQRSInput struct {
	Page   int
	Limit  int
	Search string
	Estado string
	Tipo   string
}

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

type ListPQRSResult struct {
	Items []PQRSListItem `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type ValidationError struct {
	Fields map[string]string `json:"fields"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
