package postgres

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/retry/pkg/domain"
	"github.com/Vovarama1992/retry/track-service/internal/ports"
	"github.com/sony/gobreaker"
)

type actionRepo struct {
	db      *sql.DB
	breaker *gobreaker.CircuitBreaker
}

func NewActionRepo(db *sql.DB, breaker *gobreaker.CircuitBreaker) ports.ActionRepo {
	return &actionRepo{
		db:      db,
		breaker: breaker,
	}
}

func (r *actionRepo) GetActionTypeIDByName(ctx context.Context, name string) (int64, error) {
	var id int64
	_, err := r.breaker.Execute(func() (any, error) {
		return nil, r.db.QueryRowContext(ctx,
			`SELECT id FROM action_type WHERE name = $1`, name,
		).Scan(&id)
	})
	return id, err
}

func (r *actionRepo) CreateAction(ctx context.Context, typeID int64, action domain.Action) (int64, error) {
	var id int64
	_, err := r.breaker.Execute(func() (any, error) {
		return nil, r.db.QueryRowContext(ctx,
			`INSERT INTO action (type_id, visit_id, source, ip, created_at)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id`,
			typeID, action.VisitID, action.Source, action.IPAddress, action.Timestamp,
		).Scan(&id)
	})
	return id, err
}

func (r *actionRepo) GetAllActions(ctx context.Context) ([]domain.Action, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, type_id, visit_id, source, ip, created_at FROM action`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.Action
	for rows.Next() {
		var a domain.Action
		if err := rows.Scan(&a.ID, &a.ActionTypeID, &a.VisitID, &a.Source, &a.IPAddress, &a.Timestamp); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *actionRepo) GetActionsByType(ctx context.Context, typeID int64) ([]domain.Action, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, type_id, visit_id, source, ip, created_at FROM action WHERE type_id = $1`,
		typeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.Action
	for rows.Next() {
		var a domain.Action
		if err := rows.Scan(&a.ID, &a.ActionTypeID, &a.VisitID, &a.Source, &a.IPAddress, &a.Timestamp); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *actionRepo) GetActionsByVisitID(ctx context.Context, visitID string) ([]domain.Action, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, type_id, visit_id, source, ip, created_at FROM action WHERE visit_id = $1`,
		visitID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []domain.Action
	for rows.Next() {
		var a domain.Action
		if err := rows.Scan(&a.ID, &a.ActionTypeID, &a.VisitID, &a.Source, &a.IPAddress, &a.Timestamp); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}
