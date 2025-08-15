package infra

import (
	"context"
	"database/sql"

	"github.com/Vovarama1992/retry/pkg/apperror"
	"github.com/Vovarama1992/retry/pkg/domain"
	"github.com/Vovarama1992/retry/track-service/internal/session/ports"
	"github.com/lib/pq"
	"github.com/sony/gobreaker"
)

type sessionRepo struct {
	db      *sql.DB
	breaker *gobreaker.CircuitBreaker
}

func NewSessionRepo(db *sql.DB, breaker *gobreaker.CircuitBreaker) ports.SessionRepo {
	return &sessionRepo{
		db:      db,
		breaker: breaker,
	}
}

func (r *sessionRepo) GetActionsGroupedBySessionID(ctx context.Context, limit, offset int) (map[string][]domain.Action, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		// Сначала получаем ограниченный список session_id (только непустые)
		sessionsRows, err := r.db.QueryContext(ctx, `
			SELECT DISTINCT session_id
			FROM actions
			WHERE session_id IS NOT NULL AND session_id <> ''
			ORDER BY session_id
			LIMIT $1 OFFSET $2
		`, limit, offset)
		if err != nil {
			return nil, apperror.Internal("failed to query session ids")
		}
		defer sessionsRows.Close()

		var sessionIDs []string
		for sessionsRows.Next() {
			var id string
			if err := sessionsRows.Scan(&id); err != nil {
				return nil, apperror.Internal("failed to scan session id")
			}
			sessionIDs = append(sessionIDs, id)
		}
		if len(sessionIDs) == 0 {
			return nil, apperror.NotFound("no actions found")
		}

		// Теперь получаем действия только для этих session_id
		rows, err := r.db.QueryContext(ctx, `
			SELECT id, action_type_id, visit_id, session_id, source, ip_address, timestamp, meta
			FROM actions
			WHERE session_id = ANY($1)
			ORDER BY session_id, timestamp, id
		`, pq.Array(sessionIDs))
		if err != nil {
			return nil, apperror.Internal("failed to query actions")
		}
		defer rows.Close()

		result := make(map[string][]domain.Action)
		for rows.Next() {
			var a domain.Action
			if err := rows.Scan(
				&a.ID,
				&a.ActionTypeID,
				&a.VisitID,
				&a.SessionID,
				&a.Source,
				&a.IPAddress,
				&a.Timestamp,
				&a.Meta,
			); err != nil {
				return nil, apperror.Internal("failed to scan action row")
			}
			result[a.SessionID] = append(result[a.SessionID], a)
		}

		return result, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(map[string][]domain.Action), nil
}

func (r *sessionRepo) GetSessionCountByVisitID(ctx context.Context) (map[string]int, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		rows, err := r.db.QueryContext(ctx, `
			SELECT visit_id, COUNT(DISTINCT session_id)
			FROM actions
			GROUP BY visit_id
		`)
		if err != nil {
			return nil, apperror.Internal("failed to query session counts")
		}
		defer rows.Close()

		result := make(map[string]int)
		for rows.Next() {
			var visitID string
			var count int
			if err := rows.Scan(&visitID, &count); err != nil {
				return nil, apperror.Internal("failed to scan session count row")
			}
			result[visitID] = count
		}

		if len(result) == 0 {
			return nil, apperror.NotFound("no sessions found for any visit_id")
		}

		return result, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(map[string]int), nil
}

func (r *sessionRepo) GetSessionStats(ctx context.Context) (domain.SessionStats, error) {
	res, err := r.breaker.Execute(func() (any, error) {
		var stats domain.SessionStats
		err := r.db.QueryRowContext(ctx, `
WITH sessions AS (
    SELECT
        session_id,
        MIN(timestamp) AS start_time,
        MAX(timestamp) AS end_time,
        COUNT(*) AS action_count
    FROM actions
    WHERE session_id IS NOT NULL AND session_id <> ''
    GROUP BY session_id
)
SELECT
    COUNT(*) AS total_sessions,
    COALESCE(SUM(action_count), 0) AS total_actions,
    COALESCE(AVG(EXTRACT(EPOCH FROM (end_time - start_time))), 0) AS avg_duration_seconds,
    COALESCE(AVG(action_count), 0) AS avg_actions_per_session,
    COALESCE(MAX(EXTRACT(EPOCH FROM (end_time - start_time))), 0) AS max_duration_seconds,
    COALESCE(MAX(action_count), 0) AS max_actions_per_session,
    MIN(start_time) AS first_session_at,
    MAX(end_time) AS last_session_at
FROM sessions
		`).Scan(
			&stats.TotalSessions,
			&stats.TotalActions,
			&stats.AvgDurationSeconds,
			&stats.AvgActionsPerSession,
			&stats.MaxDurationSeconds,
			&stats.MaxActionsPerSession,
			&stats.FirstSessionAt,
			&stats.LastSessionAt,
		)
		if err != nil {
			return nil, apperror.Internal("failed to fetch session stats")
		}

		if stats.TotalSessions == 0 {
			return nil, apperror.NotFound("no session stats available")
		}

		return stats, nil
	})
	if err != nil {
		return domain.SessionStats{}, err
	}
	return res.(domain.SessionStats), nil
}
