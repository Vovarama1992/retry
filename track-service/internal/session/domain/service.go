package service

import (
	"context"
	"fmt"
	"sort"
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

	out := make(map[string]summary.VisitBlock, len(visitIDs))

	for _, vID := range visitIDs {
		actions := byVisit[vID]
		if len(actions) == 0 {
			continue
		}

		// собираем временную мапу сессий
		sessionsMap := make(map[string]*summary.SessionBlock, 8)
		var visitLast time.Time

		for _, a := range actions {
			line := utils.HumanActionLine(a.Timestamp, a.ActionTypeName, a.Meta, nil)

			sb, ok := sessionsMap[a.SessionID]
			if !ok {
				sb = &summary.SessionBlock{
					SessionID:    a.SessionID,
					Actions:      make([]string, 0, 8),
					LastActionAt: time.Time{},
				}
				sessionsMap[a.SessionID] = sb
			}

			sb.Actions = append(sb.Actions, line)
			if a.Timestamp.After(sb.LastActionAt) {
				sb.LastActionAt = a.Timestamp
			}
			if a.Timestamp.After(visitLast) {
				visitLast = a.Timestamp
			}
		}

		// превращаем мапу в срез
		sessions := make([]summary.SessionBlock, 0, len(sessionsMap))
		for _, sb := range sessionsMap {
			sessions = append(sessions, *sb)
		}

		// сортировка по последнему действию (новые сверху)
		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].LastActionAt.After(sessions[j].LastActionAt)
		})

		key := vID
		if actions[0].IPAddress != "" {
			key = fmt.Sprintf("%s [%s]", vID, actions[0].IPAddress)
		}

		out[key] = summary.VisitBlock{
			Sessions:     sessions,
			LastActionAt: visitLast,
		}
	}

	return visitIDs, out, nil
}

func (s *sessionService) GetSessionCountByVisitID(ctx context.Context) (map[string]int, error) {
	return s.repo.GetSessionCountByVisitID(ctx)
}

func (s *sessionService) GetSessionStats(ctx context.Context) (domain.SessionStats, error) {
	return s.repo.GetSessionStats(ctx)
}
