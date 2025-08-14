package domain

import (
	"time"
)

type SessionStats struct {
	TotalSessions        int       `json:"total_sessions"`
	TotalActions         int       `json:"total_actions"`
	AvgDurationSeconds   float64   `json:"avg_duration_seconds"`
	AvgActionsPerSession float64   `json:"avg_actions_per_session"`
	MaxDurationSeconds   float64   `json:"max_duration_seconds"`
	MaxActionsPerSession int       `json:"max_actions_per_session"`
	FirstSessionAt       time.Time `json:"first_session_at"`
	LastSessionAt        time.Time `json:"last_session_at"`
}
