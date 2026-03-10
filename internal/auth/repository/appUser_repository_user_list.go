package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
)

func (r *AppUserInstance) GetByDisplayName(ctx context.Context, displayName string, limit, offset int) (*dto.ListAppUsersResponse, error) {
	limit, offset = normalizePagination(limit, offset, 10, 100)

	query := `
		SELECT id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
		FROM sicou.app_user
		WHERE display_name ILIKE $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, "%"+displayName+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dto.AppUser, 0, limit)
	for rows.Next() {
		var user dto.AppUser
		err := scanAppUser(rows, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sicou.app_user
		WHERE display_name ILIKE $1
	`, "%"+displayName+"%").Scan(&total); err != nil {
		return nil, err
	}

	var activeCount int64
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sicou.app_user
		WHERE display_name ILIKE $1
			AND is_active = true
	`, "%"+displayName+"%").Scan(&activeCount); err != nil {
		return nil, err
	}

	return &dto.ListAppUsersResponse{
		Items:       users,
		Total:       total,
		ActiveCount: activeCount,
		Limit:       limit,
		Offset:      offset,
		HasNext:     hasNext(offset, len(users), total),
	}, nil
}

func (r *AppUserInstance) List(ctx context.Context, limit, offset int) (*dto.ListAppUsersResponse, error) {
	limit, offset = normalizePagination(limit, offset, 10, 100)

	query := `
		SELECT id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
		FROM sicou.app_user
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dto.AppUser, 0, limit)
	for rows.Next() {
		var user dto.AppUser
		err := scanAppUser(rows, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sicou.app_user
	`).Scan(&total); err != nil {
		return nil, err
	}

	var activeCount int64
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sicou.app_user
		WHERE is_active = true
	`).Scan(&activeCount); err != nil {
		return nil, err
	}

	return &dto.ListAppUsersResponse{
		Items:       users,
		Total:       total,
		ActiveCount: activeCount,
		Limit:       limit,
		Offset:      offset,
		HasNext:     hasNext(offset, len(users), total),
	}, nil
}

func (r *AppUserInstance) ListWithFilters(ctx context.Context, req dto.ListAppUsersFilterRequest) (*dto.ListAppUsersResponse, error) {
	limit, offset := normalizePagination(req.Limit, req.Offset, 10, 100)

	search := strings.TrimSpace(req.Query)
	role := dto.NormalizeRoleCode(req.Role)

	conditions := make([]string, 0, 2)
	args := make([]any, 0, 6)
	argPos := 1

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(`
			(
				u.display_name ILIKE '%%' || $%d || '%%'
				OR u.email::text ILIKE '%%' || $%d || '%%'
			)
		`, argPos, argPos))
		args = append(args, search)
		argPos++
	}

	if role != "" {
		conditions = append(conditions, fmt.Sprintf(`
			EXISTS (
				SELECT 1
				FROM sicou.user_role ur
				INNER JOIN sicou.role ro ON ro.id = ur.role_id
				WHERE ur.user_id = u.id
					AND upper(ro.code) = $%d
			)
		`, argPos))
		args = append(args, role)
		argPos++
	}

	whereClause := "TRUE"
	if len(conditions) > 0 {
		whereClause = strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT id, email, display_name, password_hash, is_active, last_access_at, created_at, updated_at
		FROM sicou.app_user u
		WHERE %s
		ORDER BY u.display_name ASC, u.id ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	queryArgs := append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dto.AppUser, 0, limit)
	for rows.Next() {
		var user dto.AppUser
		if err := scanAppUser(rows, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sicou.app_user u
		WHERE %s
	`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	activeCountQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sicou.app_user u
		WHERE (%s) AND u.is_active = true
	`, whereClause)

	var activeCount int64
	if err := r.db.QueryRow(ctx, activeCountQuery, args...).Scan(&activeCount); err != nil {
		return nil, err
	}

	return &dto.ListAppUsersResponse{
		Items:       users,
		Total:       total,
		ActiveCount: activeCount,
		Limit:       limit,
		Offset:      offset,
		HasNext:     hasNext(offset, len(users), total),
	}, nil
}

func (r *AppUserInstance) ListOptions(ctx context.Context, req dto.ListAppUserOptionsRequest) (*dto.ListAppUserOptionsResponse, error) {
	limit, offset := normalizePagination(req.Limit, req.Offset, 20, 100)

	roles := dto.NormalizeRoleCodes(req.Roles)
	if len(roles) == 0 {
		return nil, fmt.Errorf("roles are required")
	}

	conditions := make([]string, 0, 3)
	args := make([]any, 0, 5)
	argPos := 1

	rolesArgPos := argPos
	args = append(args, roles)
	argPos++

	conditions = append(conditions, fmt.Sprintf(`
		EXISTS (
			SELECT 1
			FROM sicou.user_role ur
			INNER JOIN sicou.role ro ON ro.id = ur.role_id
			WHERE ur.user_id = u.id
				AND upper(ro.code) = ANY($%d)
		)
	`, rolesArgPos))

	conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", argPos))
	args = append(args, req.IsActive)
	argPos++

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchArgPos := argPos
		args = append(args, search)
		argPos++

		conditions = append(conditions, fmt.Sprintf(`
			(
				u.display_name ILIKE '%%' || $%d || '%%'
				OR u.email::text ILIKE '%%' || $%d || '%%'
				OR EXISTS (
					SELECT 1
					FROM sicou.citizen c
					WHERE c.deleted_at IS NULL
						AND c.email IS NOT NULL
						AND lower(c.email::text) = lower(u.email::text)
						AND c.document_number ILIKE '%%' || $%d || '%%'
				)
			)
		`, searchArgPos, searchArgPos, searchArgPos))
	}

	whereClause := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.display_name,
			u.email::text,
			(
				SELECT ro.code
				FROM sicou.user_role ur
				INNER JOIN sicou.role ro ON ro.id = ur.role_id
				WHERE ur.user_id = u.id
					AND upper(ro.code) = ANY($%d)
				ORDER BY ro.code ASC
				LIMIT 1
			) AS role
		FROM sicou.app_user u
		WHERE %s
		ORDER BY u.display_name ASC, u.id ASC
		LIMIT $%d OFFSET $%d
	`, rolesArgPos, whereClause, argPos, argPos+1)

	queryArgs := append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dto.AppUserOptionItem, 0, limit)
	for rows.Next() {
		var item dto.AppUserOptionItem
		if err := rows.Scan(&item.ID, &item.DisplayName, &item.Email, &item.Role); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sicou.app_user u
		WHERE %s
	`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	return &dto.ListAppUserOptionsResponse{
		Items:   items,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasNext: hasNext(offset, len(items), total),
	}, nil
}

func normalizePagination(limit, offset, defaultLimit, maxLimit int) (int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func hasNext(offset, currentCount int, total int64) bool {
	return int64(offset)+int64(currentCount) < total
}
