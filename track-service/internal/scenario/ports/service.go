package ports

import (
	"context"
	"time"

	"github.com/Vovarama1992/retry/track-service/internal/scenario/models"
)

type ScenarioService interface {
	// Возвращает агрегированную статистику сценария "Получить доступ"
	// включая вложенный блок "Перейти к оплате"
	GetScenarioGetAccess(ctx context.Context, limit, offset int, since time.Time) (models.ScenarioGetAccessSummary, error)
}
