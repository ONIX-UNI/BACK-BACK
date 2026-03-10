package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
)

func (r *AppUserInstance) CreateUserAuditLog(ctx context.Context, entry dto.UserAuditLogEntry) error {
	metadata := entry.Metadata
	if metadata.Before == nil {
		metadata.Before = map[string]any{}
	}
	if metadata.After == nil {
		metadata.After = map[string]any{}
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO sicou.user_audit_log (
			actor_user_id,
			actor_display_name,
			actor_email,
			action,
			target_user_id,
			target_display_name,
			target_email,
			detail,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb)
	`,
		nullableUUID(entry.ActorUserID),
		nullableText(entry.ActorName),
		nullableText(entry.ActorEmail),
		dto.NormalizeAuditAction(entry.Action),
		nullableUUID(entry.TargetUserID),
		strings.TrimSpace(entry.TargetUserName),
		strings.TrimSpace(entry.TargetUserEmail),
		strings.TrimSpace(entry.Detail),
		string(metadataJSON),
	)

	return err
}

func (r *AppUserInstance) ListUserAuditLog(ctx context.Context, req dto.ListUserAuditLogRequest) (*dto.ListUserAuditLogResponse, error) {
	limit := req.Limit
	offset := req.Offset
	limit, offset = normalizePagination(limit, offset, 20, 100)

	conditions := make([]string, 0, 7)
	args := make([]any, 0, 9)
	argPos := 1
	conditions = append(conditions, "1=1")

	if req.From != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, req.From.UTC())
		argPos++
	}
	if req.To != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, req.To.UTC())
		argPos++
	}

	if action := dto.NormalizeAuditAction(req.Action); action != "" {
		conditions = append(conditions, fmt.Sprintf("action = $%d", argPos))
		args = append(args, action)
		argPos++
	}

	if req.ActorUserID != nil {
		conditions = append(conditions, fmt.Sprintf("actor_user_id = $%d", argPos))
		args = append(args, *req.ActorUserID)
		argPos++
	}

	if req.TargetUserID != nil {
		conditions = append(conditions, fmt.Sprintf("target_user_id = $%d", argPos))
		args = append(args, *req.TargetUserID)
		argPos++
	}

	if search := strings.TrimSpace(req.Search); search != "" {
		conditions = append(conditions, fmt.Sprintf(`(
			COALESCE(actor_display_name, '') ILIKE '%%' || $%d || '%%'
			OR COALESCE(actor_email, '') ILIKE '%%' || $%d || '%%'
			OR COALESCE(target_display_name, '') ILIKE '%%' || $%d || '%%'
			OR COALESCE(target_email, '') ILIKE '%%' || $%d || '%%'
			OR COALESCE(action, '') ILIKE '%%' || $%d || '%%'
			OR COALESCE(detail, '') ILIKE '%%' || $%d || '%%'
		)`, argPos, argPos, argPos, argPos, argPos, argPos))
		args = append(args, search)
		argPos++
	}

	whereClause := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT
			id::text,
			created_at,
			COALESCE(actor_user_id::text, ''),
			COALESCE(actor_display_name, ''),
			COALESCE(actor_email, ''),
			action,
			COALESCE(target_user_id::text, ''),
			target_display_name,
			target_email,
			detail,
			metadata
		FROM sicou.user_audit_log
		WHERE %s
		ORDER BY created_at DESC, id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	queryArgs := append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dto.UserAuditItem, 0, limit)
	for rows.Next() {
		var (
			item          dto.UserAuditItem
			actorID       string
			targetID      string
			metadataBytes []byte
		)

		if err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&actorID,
			&item.Actor.Name,
			&item.Actor.Email,
			&item.Action,
			&targetID,
			&item.TargetUser.Name,
			&item.TargetUser.Email,
			&item.Detail,
			&metadataBytes,
		); err != nil {
			return nil, err
		}

		item.ActionLabel = userAuditActionLabel(item.Action)
		item.Actor.ID = strings.TrimSpace(actorID)
		item.TargetUser.ID = strings.TrimSpace(targetID)

		if len(metadataBytes) > 0 {
			_ = json.Unmarshal(metadataBytes, &item.Metadata)
		}
		if item.Metadata.Before == nil {
			item.Metadata.Before = map[string]any{}
		}
		if item.Metadata.After == nil {
			item.Metadata.After = map[string]any{}
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sicou.user_audit_log
		WHERE %s
	`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	return &dto.ListUserAuditLogResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return nil
	}
	return *id
}

func nullableText(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func userAuditActionLabel(action string) string {
	switch dto.NormalizeAuditAction(action) {
	case dto.UserAuditActionCreated:
		return "Usuario creado"
	case dto.UserAuditActionUpdated:
		return "Usuario actualizado"
	case dto.UserAuditActionDeleted:
		return "Usuario eliminado"
	case dto.UserAuditActionActivated:
		return "Usuario activado"
	case dto.UserAuditActionDeactivated:
		return "Usuario desactivado"
	case dto.UserAuditActionRoleChanged:
		return "Rol cambiado"
	case dto.UserAuditActionPasswordReset:
		return "Contrasena restablecida"
	default:
		return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(action)), "_", " ")
	}
}
