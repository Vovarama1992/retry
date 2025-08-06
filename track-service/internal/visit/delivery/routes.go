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

	// === PUBLIC ===

	// @Summary Зафиксировать визит
	// @Description Создаёт новое действие типа "visit"
	// @Accept json
	// @Produce json
	// @Param visit body VisitRequestDTO true "Информация о визите"
	// @Success 201 {object} map[string]string
	// @Failure 400,500 {string} string
	// @Router /track/visit [post]
	r.With(
		withRecover,
		withRateLimit(10, time.Minute),
	).Post("/track/visit", handler.TrackVisit)

	// @Summary Получить все визиты
	// @Produce json
	// @Success 200 {array} domain.Action
	// @Failure 500 {string} string
	// @Router /track/visit/all [get]
	r.With(
		withRecover,
		withRateLimit(5, time.Minute),
	).Get("/track/visit/all", handler.GetAllVisits)

	// @Summary Статистика по источникам визитов
	// @Produce json
	// @Success 200 {array} visit.VisitSourceStat
	// @Failure 500 {string} string
	// @Router /track/stats/by-source [get]
	r.With(
		withRecover,
		withRateLimit(5, time.Minute),
	).Get("/track/stats/by-source", handler.GetStatsBySource)
}
