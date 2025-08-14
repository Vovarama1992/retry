package visithttp

import "time"

type VisitRequestDTO struct {
	VisitID   string    `json:"visit_id" validate:"required"`
	SessionID string    `json:"session_id" validate:"required"`
	Source    string    `json:"source" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
}
