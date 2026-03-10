package handlers

import (
	"strconv"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/service"
	"github.com/gofiber/fiber/v2"
)

type SServiceTypeHandler struct {
	service *service.SServiceTypeService
}

func NewServiceTypeHandler(service *service.SServiceTypeService) *SServiceTypeHandler {
	return &SServiceTypeHandler{service: service}
}

func (h *SServiceTypeHandler) Create(c *fiber.Ctx) error {
	var req dto.ServiceTypeCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
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

func (h *SServiceTypeHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id64, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	id := int16(id64)

	var req dto.ServiceTypeUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *SServiceTypeHandler) List(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	result, err := h.service.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   result,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *SServiceTypeHandler) GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id64, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	id := int16(id64)

	result, err := h.service.GetById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *SServiceTypeHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id64, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	id := int16(id64)

	result, err := h.service.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
