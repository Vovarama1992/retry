package service

import (
	"context"
	"fmt"
	"time"

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
	// тянем визиты с экшнами
	visitIDs, byVisit, err := s.repo.GetActionsGroupedByVisitID(ctx, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	out := make(map[string]summary.VisitBlock)

	for _, vID := range visitIDs {
		actions := byVisit[vID]
		vb := summary.VisitBlock{Sessions: make(map[string][]string)}

		var last time.Time
		for _, a := range actions {
			line := utils.HumanActionLine(a.Timestamp, a.ActionTypeName, a.Meta, nil)
			vb.Sessions[a.SessionID] = append(vb.Sessions[a.SessionID], line)

			if a.Timestamp.After(last) {
				last = a.Timestamp
			}
		}

		vb.LastActionAt = last
		key := vID
		if len(actions) > 0 && actions[0].IPAddress != "" {
			key = fmt.Sprintf("%s [%s]", vID, actions[0].IPAddress)
		}
		out[key] = vb
	}

	return visitIDs, out, nil
}

func (s *sessionService) GetSessionCountByVisitID(ctx context.Context) (map[string]int, error) {
	return s.repo.GetSessionCountByVisitID(ctx)
}

func (s *sessionService) GetSessionStats(ctx context.Context) (domain.SessionStats, error) {
	return s.repo.GetSessionStats(ctx)
}
