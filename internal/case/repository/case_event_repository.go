package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SCaseEvent struct {
	db *pgxpool.Pool
}

func NewCaseEventRepository(db *pgxpool.Pool) *SCaseEvent {
	return &SCaseEvent{db: db}
}

func (r *SCaseEvent) Create(ctx context.Context, req dto.CreateCaseEventRequest) (*dto.CaseEvent, error) {
	query := `
		INSERT INTO sicou.case_event (
			case_id,
			event_type,
			title,
			notes,
			payload,
			created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, case_id, event_type, title, notes, payload, created_by, created_at
	`

	var event dto.CaseEvent

	err := r.db.QueryRow(ctx, query,
		req.CaseID,
		req.EventType,
		req.Title,
		req.Notes,
		req.Payload,
		req.CreatedBy,
	).Scan(
		&event.ID,
		&event.CaseID,
		&event.EventType,
		&event.Title,
		&event.Notes,
		&event.Payload,
		&event.CreatedBy,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}
func (r *SCaseEvent) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCaseEventRequest) (*dto.CaseEvent, error) {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.EventType != nil {
		setClauses = append(setClauses, fmt.Sprintf("event_type = $%d", argPos))
		args = append(args, *req.EventType)
		argPos++
	}

	if req.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, req.Title)
		argPos++
	}

	if req.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argPos))
		args = append(args, req.Notes)
		argPos++
	}

	if req.Payload != nil {
		setClauses = append(setClauses, fmt.Sprintf("payload = $%d", argPos))
		args = append(args, req.Payload)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE sicou.case_event
		SET %s
		WHERE id = $%d
		RETURNING id, case_id, event_type, title, notes, payload, created_by, created_at
	`, strings.Join(setClauses, ","), argPos)

	args = append(args, id)

	var event dto.CaseEvent

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&event.ID,
		&event.CaseID,
		&event.EventType,
		&event.Title,
		&event.Notes,
		&event.Payload,
		&event.CreatedBy,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}
func (r *SCaseEvent) List(ctx context.Context, limit, offset int) ([]dto.CaseEvent, error) {
	query := `
		SELECT id, case_id, event_type, title, notes, payload, created_by, created_at
		FROM sicou.case_event
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []dto.CaseEvent

	for rows.Next() {
		var event dto.CaseEvent

		err := rows.Scan(
			&event.ID,
			&event.CaseID,
			&event.EventType,
			&event.Title,
			&event.Notes,
			&event.Payload,
			&event.CreatedBy,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
func (r *SCaseEvent) GetById(ctx context.Context, id uuid.UUID) (*dto.CaseEvent, error) {
	query := `
		SELECT id, case_id, event_type, title, notes, payload, created_by, created_at
		FROM sicou.case_event
		WHERE id = $1
	`

	var event dto.CaseEvent

	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.CaseID,
		&event.EventType,
		&event.Title,
		&event.Notes,
		&event.Payload,
		&event.CreatedBy,
		&event.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}
func (r *SCaseEvent) Delete(ctx context.Context, id uuid.UUID) (*dto.CaseEvent, error) {
	query := `
		DELETE FROM sicou.case_event
		WHERE id = $1
		RETURNING id, case_id, event_type, title, notes, payload, created_by, created_at
	`

	var event dto.CaseEvent

	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.CaseID,
		&event.EventType,
		&event.Title,
		&event.Notes,
		&event.Payload,
		&event.CreatedBy,
		&event.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}
