package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/DuvanRozoParra/sicou/internal/auth/service"
	"github.com/gofiber/fiber/v2"
)

func (h *AppUserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	response, token, err := h.service.Login(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrUserInactive):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrTokenNotConfigured):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	h.setSessionCookie(c, token, response.ExpiresAt)
	return c.JSON(response)
}

func (h *AppUserHandler) Me(c *fiber.Ctx) error {
	response, err := h.service.Me(c.Context(), h.sessionCookieValue(c))
	if err != nil {
		return h.handleSessionError(c, err)
	}

	return c.JSON(response)
}

func (h *AppUserHandler) Refresh(c *fiber.Ctx) error {
	response, token, err := h.service.Refresh(c.Context(), h.sessionCookieValue(c))
	if err != nil {
		return h.handleSessionError(c, err)
	}

	h.setSessionCookie(c, token, response.ExpiresAt)
	return c.JSON(response)
}

func (h *AppUserHandler) Logout(c *fiber.Ctx) error {
	if err := h.service.Logout(c.Context(), h.sessionCookieValue(c)); err != nil {
		switch {
		case errors.Is(err, service.ErrTokenNotConfigured):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	h.clearSessionCookie(c)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AppUserHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.service.ForgotPassword(c.Context(), req); err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "email is required",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "If the email exists, reset instructions will be sent.",
	})
}

func (h *AppUserHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.service.ResetPassword(c.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrPasswordResetTokenInvalid),
			errors.Is(err, service.ErrPasswordResetTokenExpired),
			errors.Is(err, service.ErrUserInactive):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrTokenNotConfigured):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "password updated successfully",
	})
}

func (h *AppUserHandler) RequireSession(c *fiber.Ctx) error {
	response, err := h.service.Me(c.Context(), h.sessionCookieValue(c))
	if err != nil {
		return h.handleSessionError(c, err)
	}

	c.Locals("auth_user", response.User)
	c.Locals("auth_expires_at", response.ExpiresAt)
	return c.Next()
}

func (h *AppUserHandler) OptionalSession(c *fiber.Ctx) error {
	rawToken := h.sessionCookieValue(c)
	if strings.TrimSpace(rawToken) == "" {
		return c.Next()
	}

	response, err := h.service.Me(c.Context(), rawToken)
	if err == nil {
		c.Locals("auth_user", response.User)
		c.Locals("auth_expires_at", response.ExpiresAt)
	}

	return c.Next()
}

func (h *AppUserHandler) sessionCookieValue(c *fiber.Ctx) string {
	return strings.TrimSpace(c.Cookies(h.sessionCookieConf.Name))
}

func (h *AppUserHandler) setSessionCookie(c *fiber.Ctx, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.sessionCookieConf.Name,
		Value:    token,
		Path:     h.sessionCookieConf.Path,
		Domain:   h.sessionCookieConf.Domain,
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   h.sessionCookieConf.Secure,
		SameSite: h.sessionCookieConf.SameSite,
	})
}

func (h *AppUserHandler) clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     h.sessionCookieConf.Name,
		Value:    "",
		Path:     h.sessionCookieConf.Path,
		Domain:   h.sessionCookieConf.Domain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HTTPOnly: true,
		Secure:   h.sessionCookieConf.Secure,
		SameSite: h.sessionCookieConf.SameSite,
	})
}

func (h *AppUserHandler) handleSessionError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrSessionInvalid),
		errors.Is(err, service.ErrSessionExpired),
		errors.Is(err, service.ErrSessionRevoked),
		errors.Is(err, service.ErrInvalidInput),
		errors.Is(err, service.ErrUserInactive):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, service.ErrTokenNotConfigured):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
}

func normalizeSessionCookieConfig(in SessionCookieConfig) SessionCookieConfig {
	out := in

	if strings.TrimSpace(out.Name) == "" {
		out.Name = "sicou_session"
	}
	if strings.TrimSpace(out.Path) == "" {
		out.Path = "/"
	}
	out.SameSite = normalizeSameSite(out.SameSite)

	return out
}

func normalizeSameSite(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case strings.ToLower(fiber.CookieSameSiteStrictMode):
		return fiber.CookieSameSiteStrictMode
	case strings.ToLower(fiber.CookieSameSiteNoneMode):
		return fiber.CookieSameSiteNoneMode
	default:
		return fiber.CookieSameSiteLaxMode
	}
}
