package actionhttp

import (
	"encoding/json"
	"time"
)

type ActionRequestDTO struct {
	VisitID   string          `json:"visit_id" validate:"required"`
	SessionID string          `json:"session_id" validate:"required"`
	Type      string          `json:"type" validate:"required"`
	Source    string          `json:"source" validate:"required"`
	Timestamp time.Time       `json:"timestamp" validate:"required"`
	Meta      json.RawMessage `json:"meta,omitempty"`
}
