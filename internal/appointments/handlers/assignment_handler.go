package handlers

import (
	"errors"
	"log"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
	"github.com/DuvanRozoParra/sicou/internal/appointments/service"
	authdto "github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) AssignPreturno(c *fiber.Ctx) error {
	var req dto.AssignPreturnoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "INVALID_REQUEST_BODY",
			Message: "invalid request body",
			Detail:  err.Error(),
		})
	}

	actor, ok := authActor(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "missing authenticated session user",
		})
	}

	result, err := h.service.AssignPreturno(c.Context(), service.AssignPreturnoInput{
		PreturnoID:    c.Params("preturno_id"),
		CoordinatorID: req.CoordinatorID,
		ServiceTypeID: req.ServiceTypeID,
		Observations:  req.Observations,
		ActorUserID:   actor.ID.String(),
		ActorRoles:    actor.Roles,
	})
	if err != nil {
		var validationErr *service.ValidationError
		if errors.As(err, &validationErr) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
				Error:   "VALIDATION_ERROR",
				Message: validationErr.Error(),
				Fields:  validationErr.Fields,
			})
		}

		switch {
		case errors.Is(err, service.ErrAssignmentForbidden):
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error:   "ASSIGNMENT_FORBIDDEN",
				Message: "user does not have permission to assign preturnos",
			})
		case errors.Is(err, repository.ErrPreturnoNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "PRETURNO_NOT_FOUND",
				Message: "preturno not found",
			})
		case errors.Is(err, repository.ErrCoordinatorNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "COORDINATOR_NOT_FOUND",
				Message: "coordinator not found",
			})
		case errors.Is(err, repository.ErrServiceTypeNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "SERVICE_TYPE_NOT_FOUND",
				Message: "service type not found",
			})
		case errors.Is(err, repository.ErrPreturnoAssignmentConflict):
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "PRETURNO_ASSIGNMENT_CONFLICT",
				Message: "preturno does not allow assignment in current status",
			})
		default:
			log.Printf("form/legal-advise assignment failed: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "LEGAL_ADVISE_ASSIGNMENT_FAILED",
				Message: "failed to assign legal advise preturno",
				Detail:  err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.AssignPreturnoResponse{
		ID:                    result.ID,
		Status:                result.Status,
		AssignedCoordinatorID: result.AssignedCoordinatorID,
		ServiceTypeID:         result.ServiceTypeID,
		TimelineEvent: dto.AssignmentTimelineEvent{
			ID:        result.TimelineEvent.ID,
			Title:     result.TimelineEvent.Title,
			Detail:    result.TimelineEvent.Detail,
			CreatedAt: result.TimelineEvent.CreatedAt,
		},
	})
}

func (h *Handler) AssignmentOptions(c *fiber.Ctx) error {
	actor, ok := authActor(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "missing authenticated session user",
		})
	}

	result, err := h.service.AssignmentOptions(c.Context(), actor.Roles)
	if err != nil {
		if errors.Is(err, service.ErrAssignmentForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error:   "ASSIGNMENT_FORBIDDEN",
				Message: "user does not have permission to view assignment options",
			})
		}

		log.Printf("form/legal-advise assignment-options failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "LEGAL_ADVISE_ASSIGNMENT_OPTIONS_FAILED",
			Message: "failed to load assignment options",
			Detail:  err.Error(),
		})
	}

	coordinators := make([]dto.AssignmentCoordinatorOption, 0, len(result.Coordinators))
	for _, item := range result.Coordinators {
		coordinators = append(coordinators, dto.AssignmentCoordinatorOption{
			ID:          item.ID,
			DisplayName: item.DisplayName,
			Email:       item.Email,
		})
	}

	serviceTypes := make([]dto.AssignmentServiceTypeOption, 0, len(result.ServiceTypes))
	for _, item := range result.ServiceTypes {
		serviceTypes = append(serviceTypes, dto.AssignmentServiceTypeOption{
			ID:   item.ID,
			Code: item.Code,
			Name: item.Name,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.AssignmentOptionsResponse{
		Coordinators: coordinators,
		ServiceTypes: serviceTypes,
	})
}

type actorUser struct {
	ID    uuid.UUID
	Roles []string
}

func authActor(c *fiber.Ctx) (actorUser, bool) {
	raw := c.Locals("auth_user")

	switch value := raw.(type) {
	case authdto.LoginUserResponse:
		return actorUser{
			ID:    value.ID,
			Roles: normalizeRoles(value.Roles),
		}, true
	case *authdto.LoginUserResponse:
		if value == nil {
			return actorUser{}, false
		}
		return actorUser{
			ID:    value.ID,
			Roles: normalizeRoles(value.Roles),
		}, true
	default:
		return actorUser{}, false
	}
}

func normalizeRoles(in []string) []string {
	out := make([]string, 0, len(in))
	for _, role := range in {
		normalized := strings.ToUpper(strings.TrimSpace(role))
		if normalized == "" {
			continue
		}
		out = append(out, normalized)
	}
	return out
}
