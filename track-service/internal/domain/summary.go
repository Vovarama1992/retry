package service

import (
	"time"
)

type VisitsSummaryHTTPResponse struct {
	VisitIDs []string              `json:"visit_ids"`
	Visits   map[string]VisitBlock `json:"visits"`
}

type SessionBlock struct {
	SessionID    string    `json:"session_id"`
	Actions      []string  `json:"actions"`
	LastActionAt time.Time `json:"last_action_at"`
}

type VisitBlock struct {
	Sessions     []SessionBlock `json:"sessions"`
	LastActionAt time.Time      `json:"last_action_at"`
}
