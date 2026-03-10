package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/DuvanRozoParra/sicou/internal/appointments/service"
	"github.com/gofiber/fiber/v2"
)

type SDesktopHandler struct {
	service *service.SDesktopService
}

func NewDesktopHandler(service *service.SDesktopService) *SDesktopHandler {
	return &SDesktopHandler{service: service}
}

func (h *SDesktopHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateDesktopRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	res, err := h.service.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *SDesktopHandler) UpdateByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "code param is required"})
	}

	var req dto.UpdateDesktopRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	res, err := h.service.UpdateByCode(c.Context(), code, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *SDesktopHandler) List(c *fiber.Ctx) error {
	res, err := h.service.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch desktops"})
	}
	return c.JSON(res)
}

// ✅ NUEVO: DELETE /escritorios/:code
func (h *SDesktopHandler) DeleteByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "code param is required"})
	}

	if err := h.service.DeleteByCode(c.Context(), code); err != nil {
		// aquí puedes mapear errores FK a 409 si quieres
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to delete desktop",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"ok": true})
}