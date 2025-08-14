package ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	visit "github.com/Vovarama1992/retry/track-service/internal/visit/models"
)

type Service interface {
	// Создать действие (визит, клик и т.п.) по имени типа, вернуть ID
	TrackAction(ctx context.Context, actionTypeName string, action domain.Action) (int64, error)

	// Получить все действия (с пагинацией)
	GetAllActions(ctx context.Context, limit, offset int) ([]domain.Action, error)

	GetActionsByType(ctx context.Context, actionTypeName string) ([]domain.Action, error)

	GetActionsByVisitID(ctx context.Context, visitID string) ([]domain.Action, error)

	GetVisitStatsBySource(ctx context.Context) ([]visit.VisitSourceStat, error)

	// Получить действия, сгруппированные по visit_id (с пагинацией по ключам)
	GetActionsGroupedByVisitID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error)
}
