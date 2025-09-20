package scenariohttp

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

	// GET /track/scenario/get-access — агрегированная статистика по сценарию "Получить доступ"
	r.With(
		withRecover,
		withRateLimit(20, time.Minute),
	).Get("/track/scenario/get-access", handler.GetScenarioGetAccess)
}
