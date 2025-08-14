package visit

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/retry/pkg/apperror"
	"github.com/Vovarama1992/retry/pkg/domain"
	visit "github.com/Vovarama1992/retry/track-service/internal/visit/models"
	"github.com/sony/gobreaker"
)

type VisitRepo struct {
	db      *sql.DB
	breaker *gobreaker.CircuitBreaker
}

func NewVisitRepo(db *sql.DB, breaker *gobreaker.CircuitBreaker) *VisitRepo {
	return &VisitRepo{
		db:      db,
		breaker: breaker,
	}
}

func (r *VisitRepo) GetVisitStatsBySource(ctx context.Context) ([]visit.VisitSourceStat, error) {
	var result []visit.VisitSourceStat

	_, err := r.breaker.Execute(func() (any, error) {
		rows, err := r.db.QueryContext(ctx, `
			SELECT a.source, COUNT(*) 
			FROM actions a
			JOIN action_types t ON a.action_type_id = t.id
			WHERE t.name = 'visit'
			GROUP BY a.source
			ORDER BY COUNT(*) DESC
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var stat visit.VisitSourceStat
			if err := rows.Scan(&stat.Source, &stat.Count); err != nil {
				return nil, err
			}
			result = append(result, stat)
		}

		return nil, rows.Err()
	})

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, apperror.NotFound("no visit stats found")
	}

	return result, nil
}

func (r *VisitRepo) GetAllVisits(ctx context.Context, limit, offset int) ([]domain.Action, error) {
	var result []domain.Action

	_, err := r.breaker.Execute(func() (any, error) {
		rows, err := r.db.QueryContext(ctx, `
			SELECT a.id,
			       a.action_type_id,
			       t.name AS action_type_name,
			       a.visit_id,
			       a.source,
			       a.ip_address,
			       a.timestamp
			FROM actions a
			JOIN action_types t ON a.action_type_id = t.id
			WHERE t.name = 'visit'
			ORDER BY a.timestamp DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
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
			); err != nil {
				return nil, err
			}
			result = append(result, a)
		}

		return nil, rows.Err()
	})

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, apperror.NotFound("no visits found")
	}

	return result, nil
}
