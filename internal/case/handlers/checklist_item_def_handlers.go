package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SCheckListItemDef struct {
	service *service.SCheckListItemDefService
}

func NewCheckListItemDefHandler(service *service.SCheckListItemDefService) *SCheckListItemDef {
	return &SCheckListItemDef{service: service}
}

func (h *SCheckListItemDef) Create(c *fiber.Ctx) error {
	var req dto.ChecklistItemDefCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.service.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
func (h *SCheckListItemDef) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	var req dto.ChecklistItemDefUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
func (h *SCheckListItemDef) List(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	result, err := h.service.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
func (h *SCheckListItemDef) GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := h.service.GetById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
func (h *SCheckListItemDef) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := h.service.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
