package sessionhttp

import (
	"net/http"
	"time"

	"github.com/Vovarama1992/go-utils/httputil"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	withRecover := func(h http.Handler) http.Handler {
		return httputil.RecoverMiddleware(h)
	}

	withRateLimit := func(rps int, per time.Duration) func(http.Handler) http.Handler {
		return httputil.NewRateLimiter(rps, per)
	}

	// GET /track/action/grouped-by-session — экшны, сгруппированные по session_id
	r.With(
		withRecover,
		withRateLimit(5, time.Minute),
	).Get("/track/action/grouped-by-session", handler.GetActionsGroupedBySessionID)

	r.With(
		withRecover,
		withRateLimit(5, time.Minute),
	).Get("/track/action/grouped-by-visit-readable", handler.GetVisitsSummary)

	// GET /track/session/grouped-by-visit — кол-во сессий на каждый visit_id
	r.With(
		withRecover,
		withRateLimit(5, time.Minute),
	).Get("/track/session/grouped-by-visit", handler.GetSessionCountByVisitID)

	// GET /track/session/stats — агрегированная статистика по сессиям
	r.With(
		withRecover,
		withRateLimit(5, time.Minute),
	).Get("/track/session/stats", handler.GetSessionStats)

}
