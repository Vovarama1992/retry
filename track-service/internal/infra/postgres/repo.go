package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/Vovarama1992/retry/pkg/apperror"
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
		err := r.db.QueryRowContext(ctx,
			`SELECT id FROM action_types WHERE name = $1`, name,
		).Scan(&id)
		if err == sql.ErrNoRows {
			return nil, apperror.New("action type not found", 404)
		}
		return nil, err
	})
	return id, err
}

func (r *actionRepo) CreateAction(ctx context.Context, typeID int64, action domain.Action) (int64, error) {
	var id int64
	_, err := r.breaker.Execute(func() (any, error) {
		err := r.db.QueryRowContext(ctx,
			`INSERT INTO actions (action_type_id, visit_id, session_id, source, ip_address, timestamp, meta)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 RETURNING id`,
			typeID,
			action.VisitID,
			action.SessionID,
			action.Source,
			action.IPAddress,
			action.Timestamp,
			action.Meta,
		).Scan(&id)
		return nil, err
	})
	return id, err
}

func (r *actionRepo) GetAllActions(ctx context.Context, limit, offset int) ([]domain.Action, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		rows, err := r.db.QueryContext(ctx, `
			SELECT a.id,
			       a.action_type_id,
			       at.name AS action_type_name,
			       a.visit_id,
			       a.source,
			       a.ip_address,
			       a.timestamp
			FROM actions a
			JOIN action_types at ON at.id = a.action_type_id
			ORDER BY a.id
			LIMIT $1 OFFSET $2
		`, limit, offset)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var actions []domain.Action
		for rows.Next() {
			var a domain.Action
			if err := rows.Scan(
				&a.ID,
				&a.ActionTypeID,
				&a.ActionTypeName,
				&a.VisitID,
				&a.Source,
				&a.IPAddress,
				&a.Timestamp,
			); err != nil {
				return nil, err
			}
			actions = append(actions, a)
		}
		if len(actions) == 0 {
			return nil, apperror.New("no actions found", 404)
		}
		return actions, rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]domain.Action), nil
}

func (r *actionRepo) GetActionsByType(ctx context.Context, typeID int64) ([]domain.Action, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		rows, err := r.db.QueryContext(ctx,
			`SELECT id, action_type_id, visit_id, source, ip_address, timestamp
			 FROM actions WHERE action_type_id = $1`,
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
		if len(actions) == 0 {
			return nil, apperror.New("no actions found for this type", 404)
		}
		return actions, nil
	})
	if err != nil {
		return nil, err
	}
	return res.([]domain.Action), nil
}

func (r *actionRepo) GetActionsByVisitID(ctx context.Context, visitID string) ([]domain.Action, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		rows, err := r.db.QueryContext(ctx,
			`SELECT id, action_type_id, visit_id, source, ip_address, timestamp
			 FROM actions WHERE visit_id = $1`,
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
		if len(actions) == 0 {
			return nil, apperror.New("no actions found for this visit", 404)
		}
		return actions, nil
	})
	if err != nil {
		return nil, err
	}
	return res.([]domain.Action), nil
}

func (r *actionRepo) GetActionsGroupedByVisitID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		result := make(map[string][]domain.Action)

		// Сначала выбираем visit_id с лимитом и оффсетом
		visitRows, err := r.db.QueryContext(ctx, `
			SELECT DISTINCT visit_id
			FROM actions
			ORDER BY visit_id
			LIMIT $1 OFFSET $2
		`, limit, offset)
		if err != nil {
			return nil, err
		}
		defer visitRows.Close()

		var visitIDs []string
		for visitRows.Next() {
			var id string
			if err := visitRows.Scan(&id); err != nil {
				return nil, err
			}
			visitIDs = append(visitIDs, id)
		}

		if len(visitIDs) == 0 {
			return nil, apperror.New("no actions grouped by visit found", 404)
		}

		// Загружаем все действия только для этих visit_id
		query := `
			SELECT a.id,
			       a.action_type_id,
			       at.name AS action_type_name,
			       a.visit_id,
			       a.source,
			       a.ip_address,
			       a.timestamp,
			       a.meta
			FROM actions a
			JOIN action_types at ON at.id = a.action_type_id
			WHERE a.visit_id = ANY($1)
			ORDER BY a.visit_id, a.timestamp
		`
		rows, err := r.db.QueryContext(ctx, query, pq.Array(visitIDs))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var a domain.Action
			if err := rows.Scan(
				&a.ID,
				&a.ActionTypeID,
				&a.ActionTypeName,
				&a.VisitID,
				&a.Source,
				&a.IPAddress,
				&a.Timestamp,
				&a.Meta,
			); err != nil {
				return nil, err
			}
			result[a.VisitID] = append(result[a.VisitID], a)
		}

		return result, rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.(map[string][]domain.Action), nil
}
