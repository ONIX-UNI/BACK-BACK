package models

import (
	"context"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/google/uuid"
)

type ICaseFile interface {
	Create(ctx context.Context, req dto.CreateCaseFileRequest) (*dto.CaseFile, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCaseFileRequest) (*dto.CaseFile, error)
	List(ctx context.Context, limit, offset int) ([]dto.CaseFile, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.CaseFile, error)
	Delete(ctx context.Context, id uuid.UUID) (*dto.CaseFile, error)

	ListCases(ctx context.Context) ([]dto.CaseListItem, error)

	FolderList(ctx context.Context) ([]dto.FolderResponse, error)

	GetCitizenIDByEmail(
		ctx context.Context,
		email string,
	) (string, error)
	SaveOtp(
		ctx context.Context,
		citizenID string,
		otpHash string,
		expiresAt time.Time,
	) error
	InsertEmailOutbox(
		ctx context.Context,
		to []string,
		subject string,
		body string,
	) error

	VerifyOtp(
		ctx context.Context,
		citizenID string,
		otp string,
	) error

	ListCasesByEmail(
		ctx context.Context,
		email string,
	) ([]dto.CaseListItem, error)
	// Expediente(ctx context.Context) ()
}
