package repository

import (
	"context"

	"github.com/DuvanRozoParra/sicou/internal/appointments/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type STurnAudit struct {
	db *pgxpool.Pool
}

func NewTurnAuditRepository(db *pgxpool.Pool) *STurnAudit {
	return &STurnAudit{db: db}
}

func (r *STurnAudit) GetTimeLine(ctx context.Context, id uuid.UUID) ([]dto.TurnAudit, error) {
	const query = `
		SELECT 
			id,
			turn_id,
			event_type,
			title,
			notes,
			payload,
			actor_user_id,
			occurred_at
		FROM sicou.turn_audit
		WHERE turn_id = $1
		ORDER BY occurred_at ASC;
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeline []dto.TurnAudit

	for rows.Next() {
		var item dto.TurnAudit

		err := rows.Scan(
			&item.ID,
			&item.TurnID,
			&item.EventType,
			&item.Title,
			&item.Notes,
			&item.Payload,
			&item.ActorUserID,
			&item.OccurredAt,
		)
		if err != nil {
			return nil, err
		}

		timeline = append(timeline, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timeline, nil
}
