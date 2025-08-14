package sessionhttp

import (
	"github.com/Vovarama1992/retry/pkg/domain"
)

// SessionActionsResponse — ключ session_id, значение — список действий
type SessionActionsResponse map[string][]domain.Action

// SessionCountByVisitResponse — ключ visit_id, значение — кол-во уникальных session_id
type SessionCountByVisitResponse map[string]int

// SessionStatsResponse — агрегаты по сессиям
type SessionStatsResponse domain.SessionStats
