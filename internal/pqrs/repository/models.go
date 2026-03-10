package repository

import (
	"context"
	"time"
)

const (
	pqrsStatusDB = "RADICADA"
)

type AttachmentRecord struct {
	FileObjectID string
	OriginalName string
	MimeType     string
	SizeBytes    int64
	URL          string
}

type EmailContext struct {
	InternalRecipients []string
	NotifyCitizen      bool
	CitizenEmail       string
	CitizenName        string
}

type CreateRecord struct {
	ID          string
	FromEmail   *string
	Subject     string
	Body        string
	ReceivedAt  time.Time
	Attachments []AttachmentRecord
	Email       EmailContext
}

type CreateResult struct {
	ID        string
	Radicado  string
	EstadoDB  string
	CreatedAt time.Time
}

type ListInput struct {
	Page     int
	Limit    int
	Search   string
	Statuses []string
	Tipo     string
}

type ListItem struct {
	ID          string
	Radicado    string
	Tipo        string
	Asunto      string
	Ciudadano   string
	Correo      string
	EstadoDB    string
	Responsable string
	ReceivedAt  time.Time
}

type ListResult struct {
	Items []ListItem
	Total int64
	Page  int
	Limit int
}

type Repository interface {
	Create(ctx context.Context, in CreateRecord) (CreateResult, error)
	List(ctx context.Context, in ListInput) (ListResult, error)
}
