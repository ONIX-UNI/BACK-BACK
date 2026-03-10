package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *AppUserHandler) ListAuditLog(c *fiber.Ctx) error {
	limit := parseLimit(c.Query("limit"), 20, 100)
	offsetRaw := strings.TrimSpace(c.Query("offset"))
	offset := parseNonNegativeInt(offsetRaw, 0)
	if offsetRaw == "" {
		page := parsePositiveInt(c.Query("page"), 1)
		offset = (page - 1) * limit
	}

	var fromPtr *time.Time
	fromRaw := strings.TrimSpace(c.Query("from"))
	if fromRaw != "" {
		from, err := parseISODatetime(fromRaw)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid from datetime",
			})
		}
		fromPtr = &from
	}

	var toPtr *time.Time
	toRaw := strings.TrimSpace(c.Query("to"))
	if toRaw != "" {
		to, err := parseISODatetime(toRaw)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid to datetime",
			})
		}
		toPtr = &to
	}

	var actorUserID *uuid.UUID
	actorRaw := strings.TrimSpace(c.Query("actorUserId"))
	if actorRaw != "" {
		parsed, err := uuid.Parse(actorRaw)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid actorUserId",
			})
		}
		actorUserID = &parsed
	}

	var targetUserID *uuid.UUID
	targetRaw := strings.TrimSpace(c.Query("targetUserId"))
	if targetRaw != "" {
		parsed, err := uuid.Parse(targetRaw)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid targetUserId",
			})
		}
		targetUserID = &parsed
	}

	result, err := h.service.ListAuditLog(c.Context(), dto.ListUserAuditLogRequest{
		Limit:        limit,
		Offset:       offset,
		From:         fromPtr,
		To:           toPtr,
		Action:       dto.NormalizeAuditAction(c.Query("action")),
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		Search:       strings.TrimSpace(c.Query("search")),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(result)
}
