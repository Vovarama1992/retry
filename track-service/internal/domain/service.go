package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Vovarama1992/retry/pkg/apperror"
	"github.com/Vovarama1992/retry/pkg/domain"
	"github.com/Vovarama1992/retry/track-service/internal/ports"
	visit "github.com/Vovarama1992/retry/track-service/internal/visit/models"
	visit_ports "github.com/Vovarama1992/retry/track-service/internal/visit/ports"
)

type trackService struct {
	actionRepo   ports.ActionRepo
	visitService visit_ports.VisitService
}

func NewTrackService(actionRepo ports.ActionRepo, visitService visit_ports.VisitService) ports.Service {
	return &trackService{
		actionRepo:   actionRepo,
		visitService: visitService,
	}
}

func (s *trackService) TrackAction(ctx context.Context, actionTypeName string, action domain.Action) (int64, error) {
	typeID, err := s.actionRepo.GetActionTypeIDByName(ctx, actionTypeName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperror.BadRequest("unknown action type: " + actionTypeName)
		}
		return 0, apperror.Internal("failed to get action type: " + err.Error())
	}

	id, err := s.actionRepo.CreateAction(ctx, typeID, action)
	if err != nil {
		return 0, apperror.Internal("failed to create action: " + err.Error())
	}
	return id, nil
}

func (s *trackService) GetAllActions(ctx context.Context, limit, offset int) ([]domain.Action, error) {
	return s.actionRepo.GetAllActions(ctx, limit, offset)
}

func (s *trackService) GetActionsByType(ctx context.Context, actionTypeName string) ([]domain.Action, error) {
	typeID, err := s.actionRepo.GetActionTypeIDByName(ctx, actionTypeName)
	if err != nil {
		return nil, err
	}
	return s.actionRepo.GetActionsByType(ctx, typeID)
}

func (s *trackService) GetActionsByVisitID(ctx context.Context, visitID string) ([]domain.Action, error) {
	return s.actionRepo.GetActionsByVisitID(ctx, visitID)
}

func (s *trackService) GetVisitStatsBySource(ctx context.Context) ([]visit.VisitSourceStat, error) {
	return s.visitService.GetStatsBySource(ctx)
}

func (s *trackService) GetActionsGroupedByVisitID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error) {
	return s.actionRepo.GetActionsGroupedByVisitID(ctx, limit, offset)
}
