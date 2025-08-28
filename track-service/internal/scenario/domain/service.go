package domain

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/Vovarama1992/retry/pkg/domain"
	"github.com/Vovarama1992/retry/track-service/internal/scenario/models"
	"github.com/Vovarama1992/retry/track-service/internal/scenario/ports"
)

type scenarioService struct {
	repo ports.ScenarioRepo
}

func NewScenarioService(repo ports.ScenarioRepo) ports.ScenarioService {
	return &scenarioService{repo: repo}
}

func (s *scenarioService) GetScenarioGetAccess(ctx context.Context, limit, offset int) (models.ScenarioGetAccessSummary, error) {
	visitIDs, byVisit, err := s.repo.GetClickAccessStats(ctx, limit, offset)
	if err != nil {
		return models.ScenarioGetAccessSummary{}, err
	}

	summary := models.ScenarioGetAccessSummary{
		TotalVisits:              len(visitIDs),
		SessionIndexDistribution: make(map[int]int),
		ProceedToPayment: models.ScenarioProceedToPaymentStats{
			PaymentMethodsDistribution: make(map[string]int),
		},
	}

	for _, vID := range visitIDs {
		actions := byVisit[vID]
		if len(actions) == 0 {
			continue
		}

		// сортируем по времени
		sort.Slice(actions, func(i, j int) bool {
			return actions[i].Timestamp.Before(actions[j].Timestamp)
		})

		// ищем click_cta_bottom
		var clickedAt time.Time
		var sessionID string
		for _, a := range actions {
			if a.ActionTypeName == "click_cta_bottom" {
				clickedAt = a.Timestamp
				sessionID = a.SessionID
				break
			}
		}
		if clickedAt.IsZero() {
			continue
		}

		// считаем визит, где был клик
		summary.VisitsWithClick++
		detail := models.ScenarioGetAccessDetail{
			VisitID:      vID,
			SessionIndex: indexOfSession(actions, sessionID),
			ClickedAt:    clickedAt,
		}
		summary.Details = append(summary.Details, detail)
		summary.SessionIndexDistribution[detail.SessionIndex]++

		// вложенный сценарий "перейти к оплате"
		for _, a := range actions {
			if a.ActionTypeName == "click_proceed_to_payment" {
				method := ""

				if len(a.Meta) > 0 {
					var metaMap map[string]any
					if err := json.Unmarshal(a.Meta, &metaMap); err == nil {
						if v, ok := metaMap["name"].(string); ok {
							method = v
						}
					}
				}

				summary.ProceedToPayment.VisitsWithPayment++
				summary.ProceedToPayment.PaymentMethodsDistribution[method]++
				summary.ProceedToPayment.Details = append(summary.ProceedToPayment.Details, models.ScenarioProceedToPaymentDetail{
					VisitID:       vID,
					SessionIndex:  indexOfSession(actions, a.SessionID),
					ClickedAt:     a.Timestamp,
					PaymentMethod: method,
				})
			}
		}
	}

	// финальные проценты
	if summary.TotalVisits > 0 {
		summary.ConversionFromAll = float64(summary.VisitsWithClick) / float64(summary.TotalVisits) * 100
	}
	if summary.VisitsWithClick > 0 {
		summary.ProceedToPayment.ConversionFromClicks = float64(summary.ProceedToPayment.VisitsWithPayment) / float64(summary.VisitsWithClick) * 100
	}

	return summary, nil
}

// indexOfSession — вычисляет порядковый номер сессии внутри визита
func indexOfSession(actions []domain.Action, sessionID string) int {
	if sessionID == "" {
		return 0
	}
	seen := make(map[string]bool)
	order := 0
	for _, a := range actions {
		if !seen[a.SessionID] {
			order++
			seen[a.SessionID] = true
		}
		if a.SessionID == sessionID {
			return order
		}
	}
	return 0
}
