package models

import "github.com/Vovarama1992/retry/pkg/domain"

type SessionWithActions struct {
	SessionID string          `json:"session_id"`
	Actions   []domain.Action `json:"actions"`
}
