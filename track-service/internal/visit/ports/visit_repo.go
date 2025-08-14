package visit_ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	visit "github.com/Vovarama1992/retry/track-service/internal/visit/models"
)

type VisitRepo interface {
	// Агрегация по источникам визитов
	GetVisitStatsBySource(ctx context.Context) ([]visit.VisitSourceStat, error)

	// Получить все визиты (с пагинацией)
	GetAllVisits(ctx context.Context, limit, offset int) ([]domain.Action, error)
}
