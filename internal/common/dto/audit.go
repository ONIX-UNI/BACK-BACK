package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          int64           `db:"id" json:"id"`
	OccurredAt  time.Time       `db:"occurred_at" json:"occurred_at"`
	ActorUserID *uuid.UUID      `db:"actor_user_id" json:"actor_user_id,omitempty"`
	Action      string          `db:"action" json:"action"` // INSERT, UPDATE, DELETE
	TableName   string          `db:"table_name" json:"table_name"`
	RowID       *uuid.UUID      `db:"row_id" json:"row_id,omitempty"`
	BeforeData  json.RawMessage `db:"before_data" json:"before_data,omitempty"`
	AfterData   json.RawMessage `db:"after_data" json:"after_data,omitempty"`
	IPAddress   *string         `db:"ip_address" json:"ip_address,omitempty"` 
	UserAgent   *string         `db:"user_agent" json:"user_agent,omitempty"`
}