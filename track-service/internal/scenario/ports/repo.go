package ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
)

type ScenarioRepo interface {
	// Возвращает все действия по визитам, где есть click_cta_bottom
	// (и внутри можно анализировать click_proceed_to_payment)
	GetClickAccessStats(ctx context.Context, limit, offset int) ([]string, map[string][]domain.Action, error)
}
