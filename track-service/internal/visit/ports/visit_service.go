package visit_ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	model "github.com/Vovarama1992/retry/track-service/internal/visit/models"
)

type VisitService interface {
	GetStatsBySource(ctx context.Context) ([]model.VisitSourceStat, error)
	GetAllVisits(ctx context.Context) ([]domain.Action, error)
}
