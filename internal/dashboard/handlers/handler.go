package handlers

import (
	"log"
	"strconv"
	"strings"
	"time"

	authdto "github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/dashboard/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Overview(c *fiber.Ctx) error {
	authUser, ok := authSessionUser(c.Locals("auth_user"))
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "usuario no autenticado",
		})
	}
	if len(authUser.Roles) == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "sin permiso",
		})
	}

	loc, err := time.LoadLocation(defaultString(c.Query("tz"), "America/Bogota"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "parametro tz invalido",
		})
	}

	statsDateRaw := strings.TrimSpace(c.Query("stats_date"))
	statsDate := time.Now().In(loc)
	if statsDateRaw != "" {
		parsed, err := time.ParseInLocation("2006-01-02", statsDateRaw, loc)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "parametro stats_date invalido, use YYYY-MM-DD",
			})
		}
		statsDate = parsed
	}

	pendingLimit, err := parseLimit(c.Query("pending_limit"), 4)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "parametro pending_limit invalido",
		})
	}

	deadlinesLimit, err := parseLimit(c.Query("deadlines_limit"), 4)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "parametro deadlines_limit invalido",
		})
	}

	activityLimit, err := parseLimit(c.Query("activity_limit"), 5)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "parametro activity_limit invalido",
		})
	}

	response, err := h.service.Overview(c.Context(), statsDate, loc, pendingLimit, deadlinesLimit, activityLimit)
	if err != nil {
		log.Printf("dashboard overview failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func parseLimit(value string, defaultValue int) (int, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid limit")
	}
	return parsed, nil
}

func defaultString(value string, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func authSessionUser(raw any) (authdto.LoginUserResponse, bool) {
	switch value := raw.(type) {
	case authdto.LoginUserResponse:
		return value, value.ID != uuid.Nil
	case *authdto.LoginUserResponse:
		if value == nil {
			return authdto.LoginUserResponse{}, false
		}
		return *value, value.ID != uuid.Nil
	default:
		return authdto.LoginUserResponse{}, false
	}
}
