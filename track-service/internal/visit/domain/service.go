package visit_domain

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	model "github.com/Vovarama1992/retry/track-service/internal/visit/models"
	ports "github.com/Vovarama1992/retry/track-service/internal/visit/ports"
)

type VisitService struct {
	repo ports.VisitRepo
}

func NewVisitService(repo ports.VisitRepo) *VisitService {
	return &VisitService{repo: repo}
}

func (s *VisitService) GetStatsBySource(ctx context.Context) ([]model.VisitSourceStat, error) {
	return s.repo.GetVisitStatsBySource(ctx)
}

func (s *VisitService) GetAllVisits(ctx context.Context) ([]domain.Action, error) {
	return s.repo.GetAllVisits(ctx)
}
