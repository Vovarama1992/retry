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

func (s *scenarioService) GetScenarioGetAccess(ctx context.Context, limit, offset int, since time.Time) (models.ScenarioGetAccessSummary, error) {
	visitIDs, byVisit, totalVisits, err := s.repo.GetClickAccessStats(ctx, limit, offset, &since)
	if err != nil {
		return models.ScenarioGetAccessSummary{}, err
	}

	summary := models.ScenarioGetAccessSummary{
		TotalVisits:              totalVisits,
		SessionIndexDistribution: make(map[int]int),
		ProceedToPayment: models.ScenarioProceedToPaymentStats{
			PaymentMethodsDistribution: make(map[string]int),
		},
	}

	var durations []float64

	for _, vID := range visitIDs {
		actions := byVisit[vID]
		if len(actions) == 0 {
			continue
		}

		sort.Slice(actions, func(i, j int) bool {
			return actions[i].Timestamp.Before(actions[j].Timestamp)
		})

		// первый CTA после since
		var clickedAt time.Time
		var sessionID string
		for _, a := range actions {
			if a.ActionTypeName == "click_cta_bottom" && !a.Timestamp.Before(since) {
				clickedAt = a.Timestamp
				sessionID = a.SessionID
				break
			}
		}
		if clickedAt.IsZero() {
			continue
		}

		summary.VisitsWithClick++
		detail := models.ScenarioGetAccessDetail{
			VisitID:      vID,
			SessionIndex: indexOfSession(actions, sessionID),
			ClickedAt:    clickedAt,
		}
		summary.Details = append(summary.Details, detail)
		summary.SessionIndexDistribution[detail.SessionIndex]++

		// вложенный сценарий «перейти к оплате»
		for _, a := range actions {
			if a.ActionTypeName != "click_proceed_to_payment" {
				continue
			}

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

			det := models.ScenarioProceedToPaymentDetail{
				VisitID:       vID,
				SessionIndex:  indexOfSession(actions, a.SessionID),
				ClickedAt:     a.Timestamp,
				PaymentMethod: method,
			}

			// длительность в рамках той же сессии, если после CTA
			if a.SessionID == sessionID && a.Timestamp.After(clickedAt) {
				dm := a.Timestamp.Sub(clickedAt).Minutes()
				det.DurationMinutes = &dm
				durations = append(durations, dm)
			}

			summary.ProceedToPayment.Details = append(summary.ProceedToPayment.Details, det)
		}
	}

	if summary.TotalVisits > 0 {
		summary.ConversionFromAll = float64(summary.VisitsWithClick) / float64(summary.TotalVisits) * 100
	}
	if summary.VisitsWithClick > 0 {
		summary.ProceedToPayment.ConversionFromClicks =
			float64(summary.ProceedToPayment.VisitsWithPayment) / float64(summary.VisitsWithClick) * 100
	}

	// агрегаты по длительности: Avg и P50 (медиана)
	if n := len(durations); n > 0 {
		sort.Float64s(durations)
		sum := 0.0
		for _, v := range durations {
			sum += v
		}
		summary.ProceedToPayment.AvgMinutes = sum / float64(n)
		// p50 (медиана)
		if n%2 == 1 {
			summary.ProceedToPayment.P50Minutes = durations[n/2]
		} else {
			summary.ProceedToPayment.P50Minutes = (durations[n/2-1] + durations[n/2]) / 2
		}
	}

	return summary, nil
}

// indexOfSession — порядковый номер сессии внутри визита
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
