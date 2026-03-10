package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TurnAudit struct {
	ID     uuid.UUID `db:"id" json:"id"`
	TurnID uuid.UUID `db:"turn_id" json:"turn_id"`

	EventType string          `db:"event_type" json:"event_type"`
	Title     *string         `db:"title" json:"title,omitempty"`
	Notes     *string         `db:"notes" json:"notes,omitempty"`
	Payload   json.RawMessage `db:"payload" json:"payload,omitempty"`

	ActorUserID *uuid.UUID `db:"actor_user_id" json:"actor_user_id,omitempty"`
	OccurredAt  time.Time  `db:"occurred_at" json:"occurred_at"`
}

