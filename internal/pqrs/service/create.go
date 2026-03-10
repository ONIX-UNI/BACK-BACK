package service

import (
	"context"
	"errors"

	"github.com/DuvanRozoParra/sicou/internal/pqrs/repository"
	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, in CreatePQRSInput) (CreatePQRSResult, error) {
	if s.repo == nil {
		return CreatePQRSResult{}, errors.New("pqrs repository is not initialized")
	}

	normalized, validationErr := normalizeAndValidate(in)
	if validationErr != nil {
		return CreatePQRSResult{}, validationErr
	}

	pqrsID := uuid.NewString()
	attachments, err := s.saveAttachments(ctx, in.Attachments)
	if err != nil {
		return CreatePQRSResult{}, err
	}

	bodyPayload, err := buildBodyPayload(in, attachments)
	if err != nil {
		return CreatePQRSResult{}, err
	}

	var fromEmail *string
	if normalized.Email != "" {
		email := normalized.Email
		fromEmail = &email
	}

	createResult, err := s.repo.Create(ctx, repository.CreateRecord{
		ID:          pqrsID,
		FromEmail:   fromEmail,
		Subject:     normalized.RequestDescription,
		Body:        bodyPayload,
		ReceivedAt:  normalized.SubmittedAt,
		Attachments: attachments,
		Email: repository.EmailContext{
			InternalRecipients: internalMailboxRecipients(),
			NotifyCitizen:      normalized.AllowsElectronicResponseBool && normalized.Email != "",
			CitizenEmail:       normalized.Email,
			CitizenName:        buildCitizenFullName(normalized.FirstName, normalized.FirstLastName),
		},
	})
	if err != nil {
		s.cleanupAttachments(ctx, attachments)
		return CreatePQRSResult{}, err
	}

	return CreatePQRSResult{
		ID:        createResult.ID,
		Radicado:  createResult.Radicado,
		Estado:    mapPublicStatus(createResult.EstadoDB),
		CreatedAt: createResult.CreatedAt,
	}, nil
}
