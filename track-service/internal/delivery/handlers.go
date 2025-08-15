package actionhttp

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Vovarama1992/go-utils/logger"
	"github.com/Vovarama1992/retry/pkg/domain"
	track "github.com/Vovarama1992/retry/track-service/internal/ports"
	validator "github.com/go-playground/validator/v10"
)

type Handler struct {
	trackService track.Service
	logger       logger.Logger
	limitAll     int
	limitGrouped int
}

func NewHandler(trackService track.Service, logger logger.Logger) *Handler {
	limitAll := 50
	if v := os.Getenv("TRACK_ACTIONS_ALL_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limitAll = n
		}
	}

	limitGrouped := 30
	if v := os.Getenv("TRACK_ACTIONS_GROUPED_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limitGrouped = n
		}
	}

	return &Handler{
		trackService: trackService,
		logger:       logger,
		limitAll:     limitAll,
		limitGrouped: limitGrouped,
	}
}

var validate = validator.New()

// TrackAction фиксирует пользовательское действие.
// @Tags Actions
// @Summary Зафиксировать действие
// @Accept json
// @Produce json
// @Param action body ActionRequestDTO true "Информация о действии"
// @Success 201 {object} map[string]string
// @Failure 400,500 {string} string
// @Router /track/action [post]
func (h *Handler) TrackAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ActionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	ip := ExtractIP(r)

	action := domain.Action{
		VisitID:   req.VisitID,
		Source:    req.Source,
		Timestamp: req.Timestamp,
		IPAddress: ip,
		Meta:      req.Meta,
	}

	_, err := h.trackService.TrackAction(r.Context(), req.Type, action)
	if err != nil {
		h.logger.Log(logger.LogEntry{
			Level:   "error",
			Message: err.Error(),
			Error:   err,
			Service: "track",
			Method:  "TrackAction",
		})
		http.Error(w, "Failed to track action", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "action recorded"})
}

// GetAllActions возвращает все действия.
// @Tags Actions
// @Summary Получить все действия
// @Produce json
// @Param offset query int false "Смещение выборки (offset)"
// @Success 200 {array} domain.Action
// @Failure 500 {string} string
// @Router /track/action/all [get]
func (h *Handler) GetAllActions(w http.ResponseWriter, r *http.Request) {
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

	actions, err := h.trackService.GetAllActions(r.Context(), h.limitAll, offset)
	if err != nil {
		h.logger.Log(logger.LogEntry{
			Level:   "error",
			Message: err.Error(),
			Error:   err,
			Service: "track",
			Method:  "GetAllActions",
		})
		http.Error(w, "Failed to fetch actions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(actions)
}

// GetActionsGroupedByVisitID возвращает действия, сгруппированные по visit_id.
// @Tags Actions
// @Summary Получить действия, сгруппированные по visit_id
// @Produce json
// @Param offset query int false "Смещение выборки (offset)"
// @Success 200 {object} map[string][]domain.Action
// @Failure 500 {string} string
// @Router /track/action/grouped [get]
func (h *Handler) GetActionsGroupedByVisitID(w http.ResponseWriter, r *http.Request) {
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

	grouped, err := h.trackService.GetActionsGroupedByVisitID(r.Context(), h.limitGrouped, offset)
	if err != nil {
		h.logger.Log(logger.LogEntry{
			Level:   "error",
			Message: err.Error(),
			Error:   err,
			Service: "track",
			Method:  "GetActionsGroupedByVisitID",
		})
		http.Error(w, "Failed to fetch grouped actions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(grouped)
}
