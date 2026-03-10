package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, in CreateAsesoriaInput) (CreateAsesoriaResult, error) {
	if s.repo == nil {
		return CreateAsesoriaResult{}, errors.New("legal-advise repository is not initialized")
	}

	normalized, validationErr := normalizeAndValidate(in)
	if validationErr != nil {
		return CreateAsesoriaResult{}, validationErr
	}

	intakeID := uuid.NewString()
	attachments, err := s.saveAttachments(ctx, normalized.IdentityDocument, normalized.UtilityBill)
	if err != nil {
		return CreateAsesoriaResult{}, err
	}

	payload, err := buildPayload(normalized.CreateAsesoriaInput, attachments)
	if err != nil {
		s.cleanupAttachments(ctx, attachments)
		return CreateAsesoriaResult{}, err
	}

	var createdBy *uuid.UUID
	if actorUserID := strings.TrimSpace(normalized.ActorUserID); actorUserID != "" {
		if parsedActorID, err := uuid.Parse(actorUserID); err == nil {
			createdBy = &parsedActorID
		}
	}

	createResult, err := s.repo.Create(ctx, repository.CreateRecord{
		ID:                     intakeID,
		Payload:                payload,
		ConsultationDate:       normalized.ConsultationDateAt,
		SubmittedAt:            normalized.SubmittedAtTime,
		AcceptsDataProcessing:  normalized.AcceptsDataProcessingBool,
		AuthorizesNotification: normalized.AuthorizesNotificationBool,
		NotificationEmail:      normalized.Email,
		CitizenName:            normalized.FullName,
		HeadOfHousehold:        normalized.HeadOfHouseholdBool,
		CreatedBy:              createdBy,
		EventSource:            normalizeCreateTimelineSource(normalized.TimelineSource, createdBy != nil),
		Attachments:            attachments,
	})
	if err != nil {
		s.cleanupAttachments(ctx, attachments)
		return CreateAsesoriaResult{}, err
	}

	return CreateAsesoriaResult{
		ID:             createResult.ID,
		PreturnoNumber: createResult.PreturnoNumber,
		Status:         createResult.Status,
		CreatedAt:      createResult.CreatedAt,
	}, nil
}

func normalizeCreateTimelineSource(value string, hasActor bool) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case TimelineSourcePublicForm:
		return TimelineSourcePublicForm
	case TimelineSourceInternal:
		return TimelineSourceInternal
	}

	if hasActor {
		return TimelineSourceInternal
	}
	return TimelineSourcePublicForm
}
