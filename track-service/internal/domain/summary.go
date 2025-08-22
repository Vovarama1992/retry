package service

import (
	"time"
)

type VisitsSummaryHTTPResponse struct {
	VisitIDs []string              `json:"visit_ids"`
	Visits   map[string]VisitBlock `json:"visits"`
}

type VisitBlock struct {
	Sessions     map[string][]string `json:"sessions"`
	LastActionAt time.Time           `json:"last_action_at"`
}
