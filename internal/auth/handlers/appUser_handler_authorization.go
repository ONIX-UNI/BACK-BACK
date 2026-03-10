package handlers

import (
	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/gofiber/fiber/v2"
)

var userManagementAllowedRoles = map[string]struct{}{
	"SUPER_ADMIN":       {},
	"ADMIN_CONSULTORIO": {},
}

func (h *AppUserHandler) RequireUserManagementRole(c *fiber.Ctx) error {
	raw := c.Locals("auth_user")
	if raw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authenticated session user",
		})
	}

	roles := extractRolesFromAuthUser(raw)
	if !hasAnyAllowedRole(roles, userManagementAllowedRoles) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "user does not have permission to manage users",
		})
	}

	return c.Next()
}

func extractRolesFromAuthUser(raw any) []string {
	switch value := raw.(type) {
	case dto.LoginUserResponse:
		return value.Roles
	case *dto.LoginUserResponse:
		if value == nil {
			return nil
		}
		return value.Roles
	default:
		return nil
	}
}

func hasAnyAllowedRole(roles []string, allowed map[string]struct{}) bool {
	for _, role := range roles {
		normalized := dto.NormalizeRoleCode(role)
		if normalized == "" {
			continue
		}
		if _, ok := allowed[normalized]; ok {
			return true
		}
	}
	return false
}
