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

	// для conversionFromClicks: визиты, где есть хотя бы один CTA >= since
	visitsWithCtaSince := make(map[string]struct{})
	var durations []float64

	for _, vID := range visitIDs {
		actions := byVisit[vID]
		if len(actions) == 0 {
			continue
		}
		sort.Slice(actions, func(i, j int) bool { return actions[i].Timestamp.Before(actions[j].Timestamp) })

		// --- БЛОК 1: CTA all-time ---
		var firstCtaSess string
		var firstCtaAt time.Time
		haveAnyCTA := false
		eligibleSessions := make(map[string]bool) // session_id -> есть CTA >= since

		for _, a := range actions {
			if a.ActionTypeName != "click_cta_bottom" {
				continue
			}
			if !haveAnyCTA {
				firstCtaSess = a.SessionID
				firstCtaAt = a.Timestamp
				haveAnyCTA = true
			}
			if !a.Timestamp.Before(since) {
				eligibleSessions[a.SessionID] = true
			}
		}
		if !haveAnyCTA {
			continue
		}

		// считаем CTA (за всю историю)
		out.VisitsWithClick++
		// деталь и распределение по индексу сессии — по первой сессии, где встретился CTA
		sIdx := indexOfSession(actions, firstCtaSess)
		out.Details = append(out.Details, models.ScenarioGetAccessDetail{
			VisitID:      vID,
			SessionIndex: sIdx,
			ClickedAt:    firstCtaAt,
		})
		out.SessionIndexDistribution[sIdx]++

		// визит считается eligible для конверсии, если есть хотя бы ОДИН CTA >= since в ЛЮБОЙ сессии
		if len(eligibleSessions) > 0 {
			visitsWithCtaSince[vID] = struct{}{}
		}

		// --- БЛОК 2: Proceed только в сессиях, где есть CTA >= since ---
		if len(eligibleSessions) == 0 {
			continue // в этом визите ни одной сессии с CTA после since — ничего не считаем в proceed-блоке
		}

		for _, p := range actions {
			if p.ActionTypeName != "click_proceed_to_payment" {
				continue
			}
			// учитываем только оплаты в сессиях, где есть CTA >= since
			if !eligibleSessions[p.SessionID] {
				continue
			}

			// метод оплаты
			method := ""
			if len(p.Meta) > 0 {
				var meta map[string]any
				if err := json.Unmarshal(p.Meta, &meta); err == nil {
					if v, ok := meta["name"].(string); ok {
						method = v
					}
				}
			}

			// длительность: ищем CTA в этой сессии, который >= since и раньше оплаты; берём ближайший
			if ctaAt, ok := nearestCtaAfterSinceBefore(actions, p.SessionID, p.Timestamp, since); ok {
				dm := p.Timestamp.Sub(ctaAt).Minutes()
				durations = append(durations, dm)

				out.ProceedToPayment.Details = append(out.ProceedToPayment.Details, models.ScenarioProceedToPaymentDetail{
					VisitID:         vID,
					SessionIndex:    indexOfSession(actions, p.SessionID),
					ClickedAt:       p.Timestamp, // время клика «Перейти к оплате»
					PaymentMethod:   method,
					DurationMinutes: &dm,
				})
			} else {
				// CTA после since в этой сессии есть (по eligibleSessions), но ближнего до оплаты не нашли (редко, но оставим без длительности)
				out.ProceedToPayment.Details = append(out.ProceedToPayment.Details, models.ScenarioProceedToPaymentDetail{
					VisitID:       vID,
					SessionIndex:  indexOfSession(actions, p.SessionID),
					ClickedAt:     p.Timestamp,
					PaymentMethod: method,
				})
			}

			out.ProceedToPayment.VisitsWithPayment++
			out.ProceedToPayment.PaymentMethodsDistribution[method]++
		}
	}

	// Конверсии
	if out.TotalVisits > 0 {
		out.ConversionFromAll = float64(out.VisitsWithClick) / float64(out.TotalVisits) * 100
	}
	if denom := len(visitsWithCtaSince); denom > 0 {
		out.ProceedToPayment.ConversionFromClicks =
			float64(out.ProceedToPayment.VisitsWithPayment) / float64(denom) * 100
	}

	// Среднее и медиана по длительностям
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

// Ближайший CTA в заданной сессии: >= since, <= payAt, с максимальным timestamp
func nearestCtaAfterSinceBefore(actions []domain.Action, sessionID string, payAt time.Time, since time.Time) (time.Time, bool) {
	var best time.Time
	found := false
	for _, a := range actions {
		if a.ActionTypeName != "click_cta_bottom" || a.SessionID != sessionID {
			continue
		}
		if a.Timestamp.Before(since) || a.Timestamp.After(payAt) {
			continue
		}
		if !found || a.Timestamp.After(best) {
			best = a.Timestamp
			found = true
		}
	}
	return best, found
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
