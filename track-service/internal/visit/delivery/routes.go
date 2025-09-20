package visithttp

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

	r.With(
		withRecover,
		withRateLimit(3000, time.Minute),
	).Post("/track/visit", handler.TrackVisit)

	r.With(
		withRecover,
		withRateLimit(20, time.Minute),
	).Get("/track/visit/all", handler.GetAllVisits)

	r.With(
		withRecover,
		withRateLimit(20, time.Minute),
	).Get("/track/stats/by-source", handler.GetStatsBySource)
}
