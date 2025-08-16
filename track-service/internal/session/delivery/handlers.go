package sessionhttp

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Vovarama1992/go-utils/logger"
	"github.com/Vovarama1992/retry/pkg/apperror"
	"github.com/Vovarama1992/retry/pkg/domain"
	summary "github.com/Vovarama1992/retry/track-service/internal/domain"
	"github.com/Vovarama1992/retry/track-service/internal/session/ports"
)

var _ = domain.Action{}
var _ = summary.VisitBlock{}

type Handler struct {
	sessionService ports.SessionService
	logger         logger.Logger
	limitGrouped   int
}

func NewHandler(sessionService ports.SessionService, logger logger.Logger) *Handler {
	limitGrouped := 30
	if v := os.Getenv("TRACK_ACTIONS_GROUPED_BY_SESSION_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limitGrouped = n
		}
	}

	return &Handler{
		sessionService: sessionService,
		logger:         logger,
		limitGrouped:   limitGrouped,
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

// GetActionsGroupedBySessionID возвращает действия, сгруппированные по session_id.
// @Tags Sessions
// @Summary Получить действия, сгруппированные по session_id
// @Produce json
// @Param offset query int false "Смещение выборки (offset)"
// @Success 200 {array} struct{session_id string; actions []domain.Action}
// @Failure 404,500 {string} string
// @Router /track/action/grouped-by-session [get]
func (h *Handler) GetActionsGroupedBySessionID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	sessionIDs, grouped, err := h.sessionService.GetActionsGroupedBySessionID(r.Context(), h.limitGrouped, offset)
	if err != nil {
		writeError(w, h.logger, "session", "GetActionsGroupedBySessionID", err)
		return
	}

	// Собираем список в порядке sessionIDs
	out := make([]struct {
		SessionID string          `json:"session_id"`
		Actions   []domain.Action `json:"actions"`
	}, 0, len(sessionIDs))

	for _, id := range sessionIDs {
		out = append(out, struct {
			SessionID string          `json:"session_id"`
			Actions   []domain.Action `json:"actions"`
		}{
			SessionID: id,
			Actions:   grouped[id],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// GetSessionCountByVisitID возвращает количество сессий на каждый visit_id.
// @Tags Sessions
// @Summary Получить количество сессий на каждый visit_id
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 404,500 {string} string
// @Router /track/session/grouped-by-visit [get]
func (h *Handler) GetSessionCountByVisitID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	counts, err := h.sessionService.GetSessionCountByVisitID(r.Context())
	if err != nil {
		writeError(w, h.logger, "session", "GetSessionCountByVisitID", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(counts)
}

// GetSessionStats возвращает агрегированную статистику по сессиям.
// @Tags Sessions
// @Summary Получить статистику по сессиям
// @Produce json
// @Success 200 {object} domain.SessionStats
// @Failure 404,500 {string} string
// @Router /track/session/stats [get]
func (h *Handler) GetSessionStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.sessionService.GetSessionStats(r.Context())
	if err != nil {
		writeError(w, h.logger, "session", "GetSessionStats", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}

// GetVisitsSummary возвращает сводку: "visit_id [ip]" -> { sessions: { session_id: []string } }.
// @Tags Sessions
// @Summary Получить читабельную сводку по визитам (сессии и события внутри)
// @Produce json
// @Param offset query int false "Смещение выборки (offset) по сессиям"
// @Param limit  query int false "Размер страницы (limit) по сессиям. По умолчанию из env TRACK_ACTIONS_GROUPED_BY_SESSION_LIMIT"
// @Success 200 {object} map[string]summary.VisitBlock
// @Failure 400,404,500 {string} string
// @Router /track/action/grouped-by-session-readable [get]
func (h *Handler) GetVisitsSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeError(w, h.logger, "session", "GetVisitsSummary", apperror.BadRequest("invalid offset"))
			return
		}
		offset = n
	}

	limit := h.limitGrouped
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			writeError(w, h.logger, "session", "GetVisitsSummary", apperror.BadRequest("invalid limit"))
			return
		}
		limit = n
	}

	sessionIDs, data, err := h.sessionService.GetVisitsSummary(r.Context(), limit, offset)
	if err != nil {
		writeError(w, h.logger, "session", "GetVisitsSummary", err)
		return
	}

	// Чтобы не потерять порядок сессий, отдадим его рядом.
	resp := struct {
		SessionIDs []string                      `json:"session_ids"`
		Visits     map[string]summary.VisitBlock `json:"visits"`
	}{
		SessionIDs: sessionIDs,
		Visits:     data,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
