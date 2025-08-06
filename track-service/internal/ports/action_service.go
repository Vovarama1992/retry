package ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	visit "github.com/Vovarama1992/retry/track-service/internal/visit/models"
)

type Service interface {
	// Создать действие (визит, клик и т.п.) по имени типа, вернуть ID
	TrackAction(ctx context.Context, actionTypeName string, action domain.Action) (int64, error)

	GetAllActions(ctx context.Context) ([]domain.Action, error)

	GetActionsByType(ctx context.Context, actionTypeName string) ([]domain.Action, error)

	GetActionsByVisitID(ctx context.Context, visitID string) ([]domain.Action, error)

	GetVisitStatsBySource(ctx context.Context) ([]visit.VisitSourceStat, error)
}
