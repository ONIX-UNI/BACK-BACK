package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) AssignPreturno(ctx context.Context, in AssignPreturnoInput) (AssignPreturnoResult, error) {
	if r.db == nil {
		return AssignPreturnoResult{}, errors.New("database connection is not initialized")
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return AssignPreturnoResult{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := ensurePreturnoTables(ctx, tx); err != nil {
		return AssignPreturnoResult{}, err
	}

	var currentStatus string
	err = tx.QueryRow(ctx, `
		SELECT status
		FROM sicou.preturno
		WHERE id = $1::uuid
	`, in.PreturnoID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AssignPreturnoResult{}, ErrPreturnoNotFound
		}
		return AssignPreturnoResult{}, err
	}

	if !allowsPreturnoAssignment(currentStatus) {
		return AssignPreturnoResult{}, ErrPreturnoAssignmentConflict
	}

	var coordinatorName string
	err = tx.QueryRow(ctx, `
		SELECT u.display_name
		FROM sicou.app_user u
		INNER JOIN sicou.user_role ur ON ur.user_id = u.id
		INNER JOIN sicou.role ro ON ro.id = ur.role_id
		WHERE u.id = $1::uuid
			AND u.is_active = true
			AND upper(ro.code) = 'COORDINADOR'
		ORDER BY u.created_at ASC
		LIMIT 1
	`, in.CoordinatorID).Scan(&coordinatorName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AssignPreturnoResult{}, ErrCoordinatorNotFound
		}
		return AssignPreturnoResult{}, err
	}

	var serviceTypeID int16
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM sicou.service_type
		WHERE id = $1
			AND is_active = true
		LIMIT 1
	`, in.ServiceTypeID).Scan(&serviceTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AssignPreturnoResult{}, ErrServiceTypeNotFound
		}
		return AssignPreturnoResult{}, err
	}

	observations := strings.TrimSpace(in.Observations)

	out := AssignPreturnoResult{}
	err = tx.QueryRow(ctx, `
		UPDATE sicou.preturno
		SET
			status = $2,
			assigned_coordinator_id = $3::uuid,
			service_type_id = $4,
			assigned_at = now(),
			observations = CASE WHEN $5::text <> '' THEN $5 ELSE observations END,
			updated_by = COALESCE($6::uuid, updated_by)
		WHERE id = $1::uuid
		RETURNING
			id::text,
			status,
			assigned_coordinator_id::text,
			service_type_id
	`,
		in.PreturnoID,
		StatusAsignadoPreturno,
		in.CoordinatorID,
		serviceTypeID,
		observations,
		in.AssignedBy,
	).Scan(
		&out.ID,
		&out.Status,
		&out.AssignedCoordinatorID,
		&out.ServiceTypeID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AssignPreturnoResult{}, ErrPreturnoNotFound
		}
		return AssignPreturnoResult{}, err
	}

	eventDetail := buildAssignmentEventDetail(coordinatorName, observations)
	err = tx.QueryRow(ctx, `
		INSERT INTO sicou.preturno_timeline (
			preturno_id,
			event_type,
			title,
			detail,
			created_by,
			source
		)
		VALUES (
			$1::uuid,
			'ASSIGNMENT',
			'Asignacion a coordinador de turno',
			$2,
			$3::uuid,
			$4
		)
		RETURNING id::text, title, detail, created_at
	`, in.PreturnoID, eventDetail, in.AssignedBy, TimelineSourceInternal).Scan(
		&out.TimelineEvent.ID,
		&out.TimelineEvent.Title,
		&out.TimelineEvent.Detail,
		&out.TimelineEvent.CreatedAt,
	)
	if err != nil {
		return AssignPreturnoResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return AssignPreturnoResult{}, err
	}

	return out, nil
}

func (r *PostgresRepository) AssignmentOptions(ctx context.Context) (AssignmentOptionsResult, error) {
	if r.db == nil {
		return AssignmentOptionsResult{}, errors.New("database connection is not initialized")
	}

	coordinatorRows, err := r.db.Query(ctx, `
		SELECT DISTINCT
			u.id::text,
			u.display_name,
			COALESCE(u.email::text, '')
		FROM sicou.app_user u
		INNER JOIN sicou.user_role ur ON ur.user_id = u.id
		INNER JOIN sicou.role ro ON ro.id = ur.role_id
		WHERE u.is_active = true
			AND upper(ro.code) = 'COORDINADOR'
		ORDER BY u.display_name ASC
	`)
	if err != nil {
		return AssignmentOptionsResult{}, err
	}
	defer coordinatorRows.Close()

	result := AssignmentOptionsResult{
		Coordinators: make([]AssignmentCoordinatorOption, 0),
		ServiceTypes: make([]AssignmentServiceTypeOption, 0),
	}

	for coordinatorRows.Next() {
		var item AssignmentCoordinatorOption
		if err := coordinatorRows.Scan(&item.ID, &item.DisplayName, &item.Email); err != nil {
			return AssignmentOptionsResult{}, err
		}
		result.Coordinators = append(result.Coordinators, item)
	}
	if err := coordinatorRows.Err(); err != nil {
		return AssignmentOptionsResult{}, err
	}

	serviceTypeRows, err := r.db.Query(ctx, `
		SELECT id, code, name
		FROM sicou.service_type
		WHERE is_active = true
		ORDER BY id ASC
	`)
	if err != nil {
		return AssignmentOptionsResult{}, err
	}
	defer serviceTypeRows.Close()

	for serviceTypeRows.Next() {
		var item AssignmentServiceTypeOption
		if err := serviceTypeRows.Scan(&item.ID, &item.Code, &item.Name); err != nil {
			return AssignmentOptionsResult{}, err
		}
		result.ServiceTypes = append(result.ServiceTypes, item)
	}
	if err := serviceTypeRows.Err(); err != nil {
		return AssignmentOptionsResult{}, err
	}

	return result, nil
}

func allowsPreturnoAssignment(status string) bool {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "CERRADO", "CANCELADO", "ANULADO":
		return false
	default:
		return true
	}
}

func buildAssignmentEventDetail(coordinatorName string, observations string) string {
	name := strings.TrimSpace(coordinatorName)
	if name == "" {
		name = "Coordinador"
	}

	notes := strings.TrimSpace(observations)
	if strings.EqualFold(notes, name) {
		notes = ""
	}

	if notes == "" {
		return fmt.Sprintf("Asignacion a coordinador de turno: %s", name)
	}

	return fmt.Sprintf("Asignacion a coordinador de turno: %s. Nota: %s", name, notes)
}
