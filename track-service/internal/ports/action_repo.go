package ports

import (
	"context"

	"github.com/Vovarama1992/retry/pkg/domain"
)

type ActionRepo interface {
	// Получить ID типа экшна по имени ("visit", "click")
	GetActionTypeIDByName(ctx context.Context, name string) (int64, error)

	// Создать экшн с указанным type_id и timestamp, вернуть id
	CreateAction(ctx context.Context, actionTypeID int64, action domain.Action) (int64, error)

	// Получить все экшны (с пагинацией)
	GetAllActions(ctx context.Context, limit, offset int) ([]domain.Action, error)

	// Получить экшны по типу
	GetActionsByType(ctx context.Context, actionTypeID int64) ([]domain.Action, error)

	// Получить экшны по визит ID
	GetActionsByVisitID(ctx context.Context, visitID string) ([]domain.Action, error)

	// Получить экшны, сгруппированные по визит ID (с пагинацией по ключам)
	GetActionsGroupedByVisitID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error)
}
