package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/notifications/dto"
	"github.com/DuvanRozoParra/sicou/internal/notifications/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Hello(c *fiber.Ctx) error {

	var req dto.HelloRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("invalid body")
	}

	if req.Name == "" {
		return c.Status(400).SendString("name is required")
	}

	if err := h.service.SendHello(c.Context(), req.Name); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "event sent",
	})
}
