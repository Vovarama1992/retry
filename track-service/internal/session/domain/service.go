package service

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
	"github.com/Vovarama1992/retry/track-service/internal/session/ports"
)

type sessionService struct {
	repo ports.SessionRepo
}

func NewSessionService(repo ports.SessionRepo) ports.SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) GetActionsGroupedBySessionID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error) {
	return s.repo.GetActionsGroupedBySessionID(ctx, limit, offset)
}

func (s *sessionService) GetSessionCountByVisitID(ctx context.Context) (map[string]int, error) {
	return s.repo.GetSessionCountByVisitID(ctx)
}

func (s *sessionService) GetSessionStats(ctx context.Context) (domain.SessionStats, error) {
	return s.repo.GetSessionStats(ctx)
}
