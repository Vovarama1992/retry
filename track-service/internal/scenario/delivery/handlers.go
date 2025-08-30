package scenariohttp

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Vovarama1992/go-utils/logger"
	"github.com/Vovarama1992/retry/pkg/apperror"
	models "github.com/Vovarama1992/retry/track-service/internal/scenario/models"
	"github.com/Vovarama1992/retry/track-service/internal/scenario/ports"
)

var _ models.ScenarioGetAccessSummary

type Handler struct {
	scenarioService ports.ScenarioService
	logger          logger.Logger
	limitDefault    int
}

func NewHandler(scenarioService ports.ScenarioService, logger logger.Logger) *Handler {
	limit := 50
	if v := os.Getenv("TRACK_SCENARIO_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	return &Handler{
		scenarioService: scenarioService,
		logger:          logger,
		limitDefault:    limit,
	}
}

func writeError(w http.ResponseWriter, log logger.Logger, service, method string, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	log.Log(logger.LogEntry{
		Level:   "error",
		Message: err.Error(),
		Error:   err,
		Service: service,
		Method:  method,
	})
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// GetScenarioGetAccess возвращает агрегированную статистику сценария "Получить доступ"
// @Tags Scenarios
// @Summary Получить сценарий "Получить доступ"
// @Produce json
// @Param offset query int false "Смещение выборки (offset)"
// @Param limit  query int false "Размер выборки (по умолчанию TRACK_SCENARIO_LIMIT)"
// @Success 200 {object} models.ScenarioGetAccessSummary
// @Failure 400,404,500 {string} string
// @Router /track/scenario/get-access [get]
func (h *Handler) GetScenarioGetAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeError(w, h.logger, "scenario", "GetScenarioGetAccess", apperror.BadRequest("invalid offset"))
			return
		}
		offset = n
	}

	limit := h.limitDefault
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			writeError(w, h.logger, "scenario", "GetScenarioGetAccess", apperror.BadRequest("invalid limit"))
			return
		}
		limit = n
	}

	// since: опционально (?since=YYYY-MM-DD или RFC3339). По умолчанию 2025-08-28 UTC.
	var since time.Time
	if sv := r.URL.Query().Get("since"); sv != "" {
		if t, err := time.Parse(time.RFC3339, sv); err == nil {
			since = t.UTC()
		} else if t2, err2 := time.Parse("2006-01-02", sv); err2 == nil {
			since = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			writeError(w, h.logger, "scenario", "GetScenarioGetAccess", apperror.BadRequest("invalid since"))
			return
		}
	} else {
		since = time.Date(2025, 8, 28, 0, 0, 0, 0, time.UTC)
	}

	summary, err := h.scenarioService.GetScenarioGetAccess(r.Context(), limit, offset, since)
	if err != nil {
		writeError(w, h.logger, "scenario", "GetScenarioGetAccess", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(summary)
}
