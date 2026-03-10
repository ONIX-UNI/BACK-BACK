package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	services "github.com/DuvanRozoParra/sicou/internal/documents/service"
	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct {
	service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateDocumentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	result, err := h.service.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *DocumentHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "document not found",
		})
	}

	return c.JSON(result)
}

func (h *DocumentHandler) GetByFileID(c *fiber.Ctx) error {
	fileID := c.Query("file_id")

	result, err := h.service.GetByFileID(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *DocumentHandler) GetByCaseID(c *fiber.Ctx) error {
	caseID := c.Query("case_id")

	result, err := h.service.GetByCaseID(c.Context(), caseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	err := h.service.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *DocumentHandler) Get(c *fiber.Ctx) error {
	fileID := c.Query("file_id")
	caseID := c.Query("case_id")

	// Prioridad determinística
	if fileID != "" {
		result, err := h.service.GetByFileID(c.Context(), fileID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(result)
	}

	if caseID != "" {
		result, err := h.service.GetByCaseID(c.Context(), caseID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(result)
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "at least one filter is required (file_id or case_id)",
	})
}
