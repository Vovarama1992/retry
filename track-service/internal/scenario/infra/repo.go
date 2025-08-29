package infra

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/retry/pkg/apperror"
	"github.com/Vovarama1992/retry/pkg/domain"
	"github.com/Vovarama1992/retry/track-service/internal/scenario/ports"
	"github.com/lib/pq"
	"github.com/sony/gobreaker"
)

type scenarioRepo struct {
	db      *sql.DB
	breaker *gobreaker.CircuitBreaker
}

func NewScenarioRepo(db *sql.DB, breaker *gobreaker.CircuitBreaker) ports.ScenarioRepo {
	return &scenarioRepo{db: db, breaker: breaker}
}

func (r *scenarioRepo) GetClickAccessStats(ctx context.Context, limit, offset int) ([]string, map[string][]domain.Action, int, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		// общее число визитов
		var total int
		if err := r.db.QueryRowContext(ctx, `
			SELECT COUNT(DISTINCT visit_id)
			FROM actions
			WHERE visit_id IS NOT NULL AND visit_id <> ''
		`).Scan(&total); err != nil {
			return nil, apperror.Internal("failed to count total visits")
		}

		// 1) Визиты, где был click_cta_bottom
		visitsRows, err := r.db.QueryContext(ctx, `
			WITH grouped AS (
				SELECT
					visit_id,
					MAX(timestamp) AS last_time
				FROM actions a
				JOIN action_types t ON t.id = a.action_type_id
				WHERE visit_id IS NOT NULL AND visit_id <> ''
				  AND session_id IS NOT NULL AND session_id <> ''
				  AND t.name = 'click_cta_bottom'
				GROUP BY visit_id
			)
			SELECT visit_id
			FROM grouped
			ORDER BY last_time DESC, visit_id
			LIMIT $1 OFFSET $2
		`, limit, offset)
		if err != nil {
			return nil, apperror.Internal("failed to query scenario visit ids")
		}
		defer visitsRows.Close()

		var visitIDs []string
		for visitsRows.Next() {
			var id string
			if err := visitsRows.Scan(&id); err != nil {
				return nil, apperror.Internal("failed to scan visit id")
			}
			visitIDs = append(visitIDs, id)
		}
		if err := visitsRows.Err(); err != nil {
			return nil, apperror.Internal("failed to read scenario visit ids")
		}
		if len(visitIDs) == 0 {
			return nil, apperror.NotFound("no scenario visits found")
		}

		// 2) Все действия по этим визитам
		rows, err := r.db.QueryContext(ctx, `
			WITH selected AS (SELECT UNNEST($1::text[]) AS visit_id)
			SELECT
				a.id,
				a.action_type_id,
				COALESCE(t.name, '')      AS action_type_name,
				COALESCE(a.visit_id, '')  AS visit_id,
				COALESCE(a.session_id, '') AS session_id,
				COALESCE(a.source, '')    AS source,
				COALESCE(a.ip_address, '') AS ip_address,
				a.timestamp,
				COALESCE(a.meta, '{}'::jsonb) AS meta
			FROM actions a
			JOIN selected s ON s.visit_id = a.visit_id
			LEFT JOIN action_types t ON t.id = a.action_type_id
			WHERE a.session_id IS NOT NULL AND a.session_id <> ''
			ORDER BY a.visit_id, a.session_id, a.timestamp, a.id
		`, pq.Array(visitIDs))
		if err != nil {
			return nil, apperror.Internal("failed to query actions for scenario visits")
		}
		defer rows.Close()

		result := make(map[string][]domain.Action, len(visitIDs))
		for rows.Next() {
			var a domain.Action
			if err := rows.Scan(
				&a.ID,
				&a.ActionTypeID,
				&a.ActionTypeName,
				&a.VisitID,
				&a.SessionID,
				&a.Source,
				&a.IPAddress,
				&a.Timestamp,
				&a.Meta,
			); err != nil {
				return nil, apperror.Internal("failed to scan action row")
			}
			result[a.VisitID] = append(result[a.VisitID], a)
		}
		if err := rows.Err(); err != nil {
			return nil, apperror.Internal("failed to read scenario actions rows")
		}

		return struct {
			IDs    []string
			Result map[string][]domain.Action
			Total  int
		}{visitIDs, result, total}, nil
	})
	if err != nil {
		return nil, nil, 0, err
	}

	out := res.(struct {
		IDs    []string
		Result map[string][]domain.Action
		Total  int
	})
	return out.IDs, out.Result, out.Total, nil
}
