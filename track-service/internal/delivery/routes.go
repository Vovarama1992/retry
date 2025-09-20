package actionhttp

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

	// POST /track/action — создать новое действие
	r.With(
		withRecover,
		withRateLimit(3000, time.Minute),
	).Post("/track/action", handler.TrackAction)

	// GET /track/action/all — получить все действия
	r.With(
		withRecover,
		withRateLimit(20, time.Minute),
	).Get("/track/action/all", handler.GetAllActions)

	// GET /track/action/grouped — получить все действия, сгруппированные по visit_id
	r.With(
		withRecover,
		withRateLimit(20, time.Minute),
	).Get("/track/action/grouped", handler.GetActionsGroupedByVisitID)
}
