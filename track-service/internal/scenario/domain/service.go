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
	visitIDs, byVisit, totalVisits, err := s.repo.GetClickAccessStats(ctx, limit, offset)
	if err != nil {
		return models.ScenarioGetAccessSummary{}, err
	}

	out := models.ScenarioGetAccessSummary{
		TotalVisits:              totalVisits,
		SessionIndexDistribution: make(map[int]int),
		ProceedToPayment: models.ScenarioProceedToPaymentStats{
			PaymentMethodsDistribution: make(map[string]int),
		},
	}

	var durations []float64
	clicksSince := 0 // CTA после since — знаменатель для conversionFromClicks

	for _, vID := range visitIDs {
		actions := byVisit[vID]
		if len(actions) == 0 {
			continue
		}

		sort.Slice(actions, func(i, j int) bool { return actions[i].Timestamp.Before(actions[j].Timestamp) })

		// первый CTA (за всю историю)
		var ctaAt time.Time
		var ctaSess string
		for _, a := range actions {
			if a.ActionTypeName == "click_cta_bottom" {
				ctaAt = a.Timestamp
				ctaSess = a.SessionID
				break
			}
		}
		if ctaAt.IsZero() {
			continue
		}

		// считаем CTA всегда (all-time)
		out.VisitsWithClick++
		out.Details = append(out.Details, models.ScenarioGetAccessDetail{
			VisitID:      vID,
			SessionIndex: indexOfSession(actions, ctaSess),
			ClickedAt:    ctaAt,
		})
		out.SessionIndexDistribution[indexOfSession(actions, ctaSess)]++

		// CTA после since?
		ctaEligible := !ctaAt.Before(since)
		if ctaEligible {
			clicksSince++
		}

		// оплаты считаем ТОЛЬКО если CTA ≥ since
		if !ctaEligible {
			continue
		}

		for _, a := range actions {
			if a.ActionTypeName != "click_proceed_to_payment" {
				continue
			}

			method := ""
			if len(a.Meta) > 0 {
				var meta map[string]any
				if err := json.Unmarshal(a.Meta, &meta); err == nil {
					if v, ok := meta["name"].(string); ok {
						method = v
					}
				}
			}

			out.ProceedToPayment.VisitsWithPayment++
			out.ProceedToPayment.PaymentMethodsDistribution[method]++

			det := models.ScenarioProceedToPaymentDetail{
				VisitID:       vID,
				SessionIndex:  indexOfSession(actions, a.SessionID),
				ClickedAt:     a.Timestamp, // время клика "Перейти к оплате"
				PaymentMethod: method,
			}
			// длительность — только если та же сессия и после CTA
			if a.SessionID == ctaSess && a.Timestamp.After(ctaAt) {
				dm := a.Timestamp.Sub(ctaAt).Minutes()
				det.DurationMinutes = &dm
				durations = append(durations, dm)
			}
			out.ProceedToPayment.Details = append(out.ProceedToPayment.Details, det)
		}
	}

	// конверсии
	if out.TotalVisits > 0 {
		out.ConversionFromAll = float64(out.VisitsWithClick) / float64(out.TotalVisits) * 100
	}
	if clicksSince > 0 {
		out.ProceedToPayment.ConversionFromClicks =
			float64(out.ProceedToPayment.VisitsWithPayment) / float64(clicksSince) * 100
	}

	// среднее и медиана по длительностям
	if n := len(durations); n > 0 {
		sort.Float64s(durations)
		sum := 0.0
		for _, v := range durations {
			sum += v
		}
		out.ProceedToPayment.AvgMinutes = sum / float64(n)
		if n%2 == 1 {
			out.ProceedToPayment.P50Minutes = durations[n/2]
		} else {
			out.ProceedToPayment.P50Minutes = (durations[n/2-1] + durations[n/2]) / 2
		}
	}

	return out, nil
}

// порядковый номер сессии внутри визита
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
