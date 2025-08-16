package ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	summary "github.com/Vovarama1992/retry/track-service/internal/domain"
)

// SessionService описывает бизнес-логику по сессиям.
type SessionService interface {
	// Получить действия, сгруппированные по session_id (с пагинацией по ключам)
	GetActionsGroupedBySessionID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error)

	GetVisitsSummary(ctx context.Context, limit, offset int) (map[string]summary.VisitBlock, error)

	GetSessionCountByVisitID(ctx context.Context) (map[string]int, error)
	GetSessionStats(ctx context.Context) (domain.SessionStats, error)
}
