package ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
)

// SessionRepo описывает работу с сессиями на уровне хранилища.
type SessionRepo interface {
	// GetActionsGroupedBySessionID возвращает map[session_id][]Action (с пагинацией по session_id)
	GetActionsGroupedBySessionID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error)

	// GetSessionCountByVisitID возвращает map[visit_id] -> count уникальных session_id
	GetSessionCountByVisitID(ctx context.Context) (map[string]int, error)

	// GetSessionStats возвращает агрегаты по сессиям
	GetSessionStats(ctx context.Context) (domain.SessionStats, error)
}
