package service

import (
	"context"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
)

func (s *AppUserService) ListAuditLog(ctx context.Context, req dto.ListUserAuditLogRequest) (*dto.ListUserAuditLogResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	req.Action = dto.NormalizeAuditAction(req.Action)
	req.Search = strings.TrimSpace(req.Search)

	if req.From != nil && req.To != nil && req.From.After(*req.To) {
		return nil, ErrInvalidInput
	}
	if req.Action != "" && !isAllowedUserAuditAction(req.Action) {
		return nil, ErrInvalidInput
	}

	return s.repo.ListUserAuditLog(ctx, req)
}

func (s *AppUserService) createUserAuditLog(ctx context.Context, entry dto.UserAuditLogEntry) error {
	if !isAllowedUserAuditAction(dto.NormalizeAuditAction(entry.Action)) {
		return ErrInvalidInput
	}

	entry.Action = dto.NormalizeAuditAction(entry.Action)
	entry.Detail = strings.TrimSpace(entry.Detail)
	if entry.Detail == "" {
		return ErrInvalidInput
	}

	if entry.ActorUserID != nil {
		actor, err := s.repo.GetByID(ctx, *entry.ActorUserID)
		if err != nil {
			return err
		}
		if actor != nil {
			entry.ActorName = strings.TrimSpace(actor.DisplayName)
			entry.ActorEmail = strings.TrimSpace(actor.Email)
		}
	}

	if strings.TrimSpace(entry.TargetUserName) == "" || strings.TrimSpace(entry.TargetUserEmail) == "" {
		return ErrInvalidInput
	}

	return s.repo.CreateUserAuditLog(ctx, entry)
}

func buildCreateUserDetail(targetName string, targetEmail string, role string) string {
	userLabel := formatUserLabel(targetName, targetEmail)
	value := strings.TrimSpace(role)
	if value == "" {
		return "Creo el usuario " + userLabel
	}
	return "Creo el usuario " + userLabel + " con rol " + humanizeRole(value)
}

func buildUserUpdateAuditEvents(
	before dto.AppUser,
	after dto.AppUser,
	beforeState map[string]any,
	afterState map[string]any,
	beforeRole string,
	afterRole string,
	beforeIsActive bool,
	afterIsActive bool,
) []dto.UserAuditLogEntry {
	events := make([]dto.UserAuditLogEntry, 0, 3)

	if beforeRole != afterRole {
		userLabel := formatUserLabel(after.DisplayName, after.Email)
		events = append(events, dto.UserAuditLogEntry{
			Action: dto.UserAuditActionRoleChanged,
			Detail: buildRoleChangedDetail(userLabel, beforeRole, afterRole),
			Metadata: dto.UserAuditMetadata{
				Before: map[string]any{
					"role": beforeRole,
				},
				After: map[string]any{
					"role": afterRole,
				},
			},
		})
	}

	if beforeIsActive != afterIsActive {
		userLabel := formatUserLabel(after.DisplayName, after.Email)
		action := dto.UserAuditActionDeactivated
		detail := "Desactivo el usuario " + userLabel
		if afterIsActive {
			action = dto.UserAuditActionActivated
			detail = "Activo el usuario " + userLabel
		}
		events = append(events, dto.UserAuditLogEntry{
			Action: action,
			Detail: detail,
			Metadata: dto.UserAuditMetadata{
				Before: map[string]any{
					"is_active": beforeIsActive,
				},
				After: map[string]any{
					"is_active": afterIsActive,
				},
			},
		})
	}

	if strings.TrimSpace(before.DisplayName) != strings.TrimSpace(after.DisplayName) {
		userLabel := formatUserLabel(after.DisplayName, after.Email)
		events = append(events, dto.UserAuditLogEntry{
			Action: dto.UserAuditActionUpdated,
			Detail: "Actualizo los datos del usuario " + userLabel,
			Metadata: dto.UserAuditMetadata{
				Before: beforeState,
				After:  afterState,
			},
		})
	}

	return events
}

func buildRoleChangedDetail(userLabel string, beforeRole string, afterRole string) string {
	left := strings.TrimSpace(beforeRole)
	right := strings.TrimSpace(afterRole)
	if left == "" {
		left = "SIN_ROL"
	}
	if right == "" {
		right = "SIN_ROL"
	}
	return "Cambio el rol de " + userLabel + " de " + humanizeRole(left) + " a " + humanizeRole(right)
}

func userAuditState(user dto.AppUser) map[string]any {
	return map[string]any{
		"display_name": strings.TrimSpace(user.DisplayName),
		"email":        strings.TrimSpace(user.Email),
		"is_active":    user.IsActive,
		"role":         strings.TrimSpace(user.Role),
		"roles":        append([]string{}, user.Roles...),
	}
}

func pointerUUID(id uuid.UUID) *uuid.UUID {
	value := id
	return &value
}

func isAllowedUserAuditAction(action string) bool {
	switch dto.NormalizeAuditAction(action) {
	case dto.UserAuditActionCreated,
		dto.UserAuditActionUpdated,
		dto.UserAuditActionDeleted,
		dto.UserAuditActionActivated,
		dto.UserAuditActionDeactivated,
		dto.UserAuditActionRoleChanged,
		dto.UserAuditActionPasswordReset:
		return true
	default:
		return false
	}
}

func buildDeleteUserDetail(targetName string, targetEmail string) string {
	return "Elimino el usuario " + formatUserLabel(targetName, targetEmail)
}

func formatUserLabel(name string, email string) string {
	display := strings.TrimSpace(name)
	addr := strings.TrimSpace(email)
	if display == "" && addr == "" {
		return "sin identificar"
	}
	if display == "" {
		return addr
	}
	if addr == "" {
		return display
	}
	return display + " (" + addr + ")"
}

func humanizeRole(roleCode string) string {
	switch strings.ToUpper(strings.TrimSpace(roleCode)) {
	case "":
		return "Sin rol"
	case "SIN_ROL":
		return "Sin rol"
	case "SUPER_ADMIN":
		return "Super Administrador"
	case "ADMIN_CONSULTORIO":
		return "Administrador de Consultorio"
	case "SECRETARIA":
		return "Secretaria"
	case "COORDINADOR":
		return "Coordinador"
	case "ESTUDIANTE":
		return "Estudiante"
	case "DOCENTE":
		return "Docente"
	default:
		value := strings.ToUpper(strings.TrimSpace(roleCode))
		return strings.ReplaceAll(value, "_", " ")
	}
}
