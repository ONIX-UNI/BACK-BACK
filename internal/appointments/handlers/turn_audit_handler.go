package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/appointments/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type STurnAuditHandler struct {
	service *service.STurnAuditService
}

func NewTurnAuditHandler(service *service.STurnAuditService) *STurnAuditHandler {
	return &STurnAuditHandler{service: service}
}

func (h *STurnAuditHandler) GetTimeLine(c *fiber.Ctx) error {
	idParam := c.Params("id")

	turnID, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid turn id")
	}

	timeline, err := h.service.GetTimeLine(c.Context(), turnID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(timeline)
}
