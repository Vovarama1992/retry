package service

import (
	"context"
	"fmt"

	"github.com/Vovarama1992/retry/pkg/domain"
	summary "github.com/Vovarama1992/retry/track-service/internal/domain"
	"github.com/Vovarama1992/retry/track-service/internal/session/ports"
	"github.com/Vovarama1992/retry/track-service/internal/session/utils"
)

type sessionService struct {
	repo ports.SessionRepo
}

func NewSessionService(repo ports.SessionRepo) ports.SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) GetActionsGroupedBySessionID(ctx context.Context, limit, offset int) ([]string, map[string][]domain.Action, error) {
	return s.repo.GetActionsGroupedBySessionID(ctx, limit, offset)
}

func (s *sessionService) GetVisitsSummary(ctx context.Context, limit, offset int) ([]string, map[string]summary.VisitBlock, error) {
	sessionIDs, bySession, err := s.repo.GetActionsGroupedBySessionID(ctx, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	out := make(map[string]summary.VisitBlock)
	for _, sessID := range sessionIDs { // порядок сохраняем
		actions := bySession[sessID]
		for _, a := range actions {
			key := a.VisitID
			if a.IPAddress != "" {
				key = fmt.Sprintf("%s [%s]", a.VisitID, a.IPAddress)
			}
			vb := out[key]
			if vb.Sessions == nil {
				vb.Sessions = make(map[string][]string)
			}
			line := utils.HumanActionLine(a.Timestamp, a.ActionTypeName, a.Meta, nil)
			vb.Sessions[sessID] = append(vb.Sessions[sessID], line)
			out[key] = vb
		}
	}
	return sessionIDs, out, nil
}

func (s *sessionService) GetSessionCountByVisitID(ctx context.Context) (map[string]int, error) {
	return s.repo.GetSessionCountByVisitID(ctx)
}

func (s *sessionService) GetSessionStats(ctx context.Context) (domain.SessionStats, error) {
	return s.repo.GetSessionStats(ctx)
}
