package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
	"github.com/google/uuid"
)

var ErrAssignmentForbidden = errors.New("forbidden to assign preturno")

func (s *Service) AssignPreturno(ctx context.Context, in AssignPreturnoInput) (AssignPreturnoResult, error) {
	if s.repo == nil {
		return AssignPreturnoResult{}, errors.New("legal-advise repository is not initialized")
	}

	if !hasAssignmentPermission(in.ActorRoles) {
		return AssignPreturnoResult{}, ErrAssignmentForbidden
	}

	fieldsErr := make(map[string]string)

	preturnoID := strings.TrimSpace(in.PreturnoID)
	if preturnoID == "" {
		fieldsErr["preturno_id"] = "is required"
	} else if _, err := uuid.Parse(preturnoID); err != nil {
		fieldsErr["preturno_id"] = "must be a valid uuid"
	}

	coordinatorID := strings.TrimSpace(in.CoordinatorID)
	if coordinatorID == "" {
		fieldsErr["coordinator_id"] = "is required"
	} else if _, err := uuid.Parse(coordinatorID); err != nil {
		fieldsErr["coordinator_id"] = "must be a valid uuid"
	}

	serviceTypeRaw := strings.TrimSpace(in.ServiceTypeID)
	if serviceTypeRaw == "" {
		fieldsErr["service_type_id"] = "is required"
	}

	parsedServiceTypeID, parseErr := parseServiceTypeID(serviceTypeRaw)
	if serviceTypeRaw != "" && parseErr != nil {
		fieldsErr["service_type_id"] = parseErr.Error()
	}

	if len(fieldsErr) > 0 {
		return AssignPreturnoResult{}, &ValidationError{Fields: fieldsErr}
	}

	var assignedBy *uuid.UUID
	if actorID := strings.TrimSpace(in.ActorUserID); actorID != "" {
		if parsedActorID, err := uuid.Parse(actorID); err == nil {
			assignedBy = &parsedActorID
		}
	}

	repoResult, err := s.repo.AssignPreturno(ctx, repository.AssignPreturnoInput{
		PreturnoID:    preturnoID,
		CoordinatorID: coordinatorID,
		ServiceTypeID: parsedServiceTypeID,
		Observations:  strings.TrimSpace(in.Observations),
		AssignedBy:    assignedBy,
	})
	if err != nil {
		return AssignPreturnoResult{}, err
	}

	return AssignPreturnoResult{
		ID:                    repoResult.ID,
		Status:                mapAssignmentStatus(repoResult.Status),
		AssignedCoordinatorID: repoResult.AssignedCoordinatorID,
		ServiceTypeID:         strconv.FormatInt(int64(repoResult.ServiceTypeID), 10),
		TimelineEvent: AssignmentTimelineEvent{
			ID:        repoResult.TimelineEvent.ID,
			Title:     repoResult.TimelineEvent.Title,
			Detail:    repoResult.TimelineEvent.Detail,
			CreatedAt: repoResult.TimelineEvent.CreatedAt,
		},
	}, nil
}

func (s *Service) AssignmentOptions(ctx context.Context, actorRoles []string) (AssignmentOptionsResult, error) {
	if s.repo == nil {
		return AssignmentOptionsResult{}, errors.New("legal-advise repository is not initialized")
	}

	if !hasAssignmentPermission(actorRoles) {
		return AssignmentOptionsResult{}, ErrAssignmentForbidden
	}

	repoResult, err := s.repo.AssignmentOptions(ctx)
	if err != nil {
		return AssignmentOptionsResult{}, err
	}

	coordinators := make([]AssignmentCoordinatorOption, 0, len(repoResult.Coordinators))
	for _, item := range repoResult.Coordinators {
		coordinators = append(coordinators, AssignmentCoordinatorOption{
			ID:          item.ID,
			DisplayName: item.DisplayName,
			Email:       item.Email,
		})
	}

	serviceTypes := make([]AssignmentServiceTypeOption, 0, len(repoResult.ServiceTypes))
	for _, item := range repoResult.ServiceTypes {
		serviceTypes = append(serviceTypes, AssignmentServiceTypeOption{
			ID:   strconv.FormatInt(int64(item.ID), 10),
			Code: item.Code,
			Name: item.Name,
		})
	}

	return AssignmentOptionsResult{
		Coordinators: coordinators,
		ServiceTypes: serviceTypes,
	}, nil
}

func parseServiceTypeID(value string) (int16, error) {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, errors.New("must be a valid numeric identifier")
	}
	if parsed <= 0 || parsed > 32767 {
		return 0, errors.New("must be between 1 and 32767")
	}
	return int16(parsed), nil
}

func hasAssignmentPermission(roles []string) bool {
	allowedRoles := map[string]struct{}{
		"SUPER_ADMIN":       {},
		"ADMIN_CONSULTORIO": {},
		"COORDINADOR":       {},
	}

	for _, role := range roles {
		normalized := strings.ToUpper(strings.TrimSpace(role))
		if _, ok := allowedRoles[normalized]; ok {
			return true
		}
	}
	return false
}

func mapAssignmentStatus(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case repository.StatusAsignadoPreturno:
		return "asignado"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}
