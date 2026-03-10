package handlers

import (
	"errors"
	"net/url"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *AppUserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateAppUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.service.Create(c.Context(), req, h.authActorID(c))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrRoleNotFound):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrEmailAlreadyUsed):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AppUserHandler) BulkCreate(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}
	if fileHeader.Size > maxBulkUploadBytes {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": "file is too large",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to open file",
		})
	}
	defer file.Close()

	result, err := h.service.BulkCreate(c.Context(), file, h.authActorID(c))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrFileTooLarge):
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(result)
}

func (h *AppUserHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	user, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(user)
}

func (h *AppUserHandler) GetByEmail(c *fiber.Ctx) error {
	email, err := url.PathUnescape(c.Params("email"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid email path parameter",
		})
	}

	user, err := h.service.GetByEmail(c.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(user)
}

func (h *AppUserHandler) GetByDisplayName(c *fiber.Ctx) error {
	displayName := strings.TrimSpace(c.Query("display_name"))
	if displayName == "" {
		displayName = strings.TrimSpace(c.Query("search"))
	}
	limit := parseLimit(c.Query("limit"), 10, 100)
	offset := parseNonNegativeInt(c.Query("offset"), 0)

	result, err := h.service.GetByDisplayName(c.Context(), displayName, limit, offset)
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

func (h *AppUserHandler) List(c *fiber.Ctx) error {
	page, err := parseStrictPositiveInt(c.Query("page"), 1)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid page query parameter",
		})
	}

	limit, err := parseStrictLimit(c.Query("limit"), 10, 100)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit query parameter",
		})
	}

	offset := (page - 1) * limit
	rawOffset := strings.TrimSpace(c.Query("offset"))
	if rawOffset != "" {
		parsedOffset, parseErr := parseStrictNonNegativeInt(rawOffset)
		if parseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid offset query parameter",
			})
		}
		offset = parsedOffset
	}

	result, err := h.service.ListWithFilters(c.Context(), dto.ListAppUsersFilterRequest{
		Query:  strings.TrimSpace(c.Query("q")),
		Role:   dto.NormalizeRoleCode(c.Query("role")),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid role query parameter",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(result)
}

func (h *AppUserHandler) ListOptions(c *fiber.Ctx) error {
	rawRoles := strings.TrimSpace(c.Query("roles"))
	if rawRoles == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "roles is required",
		})
	}

	isActive, err := parseOptionalBool(c.Query("is_active"), true)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid is_active query parameter",
		})
	}

	limit := parseLimit(c.Query("limit"), 20, 100)
	offset := parseNonNegativeInt(c.Query("offset"), 0)

	result, err := h.service.ListOptions(c.Context(), dto.ListAppUserOptionsRequest{
		Roles:    strings.Split(rawRoles, ","),
		Search:   strings.TrimSpace(c.Query("search")),
		IsActive: isActive,
		Limit:    limit,
		Offset:   offset,
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

func (h *AppUserHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	var req dto.UpdateAppUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.service.Update(c.Context(), id, req, h.authActorID(c))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(user)
}

func (h *AppUserHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	err = h.service.Delete(c.Context(), id, h.authActorID(c))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AppUserHandler) authActorID(c *fiber.Ctx) *uuid.UUID {
	raw := c.Locals("auth_user")
	switch value := raw.(type) {
	case dto.LoginUserResponse:
		id := value.ID
		return &id
	case *dto.LoginUserResponse:
		if value == nil {
			return nil
		}
		id := value.ID
		return &id
	default:
		return nil
	}
}
