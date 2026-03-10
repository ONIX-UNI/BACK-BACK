package handlers

import (
	"strconv"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/service"
	"github.com/gofiber/fiber/v2"
)

type SLegalAreaHandler struct {
	service *service.SlegalAreaService
}

func NewLegalAreaHandler(service *service.SlegalAreaService) *SLegalAreaHandler {
	return &SLegalAreaHandler{service: service}
}

func (h *SLegalAreaHandler) Create(c *fiber.Ctx) error {

	var req dto.CreateLegalAreaRequest

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
func (h *SLegalAreaHandler) Update(c *fiber.Ctx) error {

	idParam := c.Params("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var req dto.UpdateLegalAreaRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	result, err := h.service.Update(c.Context(), int16(idInt), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
func (h *SLegalAreaHandler) List(c *fiber.Ctx) error {

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	result, err := h.service.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
func (h *SLegalAreaHandler) GetById(c *fiber.Ctx) error {

	idParam := c.Params("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	result, err := h.service.GetById(c.Context(), int16(idInt))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
func (h *SLegalAreaHandler) Delete(c *fiber.Ctx) error {

	idParam := c.Params("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	result, err := h.service.Delete(c.Context(), int16(idInt))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
