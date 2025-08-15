package visithttp

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Vovarama1992/go-utils/logger"
	"github.com/Vovarama1992/retry/pkg/apperror"
	"github.com/Vovarama1992/retry/pkg/domain"
	actionhttp "github.com/Vovarama1992/retry/track-service/internal/delivery"
	track "github.com/Vovarama1992/retry/track-service/internal/ports"
	visit "github.com/Vovarama1992/retry/track-service/internal/visit/models"
	visit_ports "github.com/Vovarama1992/retry/track-service/internal/visit/ports"
	validator "github.com/go-playground/validator/v10"
)

var _ = visit.VisitSourceStat{}

type Handler struct {
	trackService track.Service
	visitService visit_ports.VisitService
	logger       logger.Logger
	limitAll     int
}

func NewHandler(trackService track.Service, visitService visit_ports.VisitService, logger logger.Logger) *Handler {
	limitAll := 50
	if v := os.Getenv("TRACK_VISITS_ALL_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limitAll = n
		}
	}

	return &Handler{
		trackService: trackService,
		visitService: visitService,
		logger:       logger,
		limitAll:     limitAll,
	}
}

var validate = validator.New()

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

// TrackVisit фиксирует визит.
// @Tags Visits
// @Summary Зафиксировать визит
// @Description Создаёт новое действие типа "visit"
// @Accept json
// @Produce json
// @Param visit body VisitRequestDTO true "Информация о визите"
// @Success 201 {object} map[string]string
// @Failure 400,404,500 {string} string
// @Router /track/visit [post]
func (h *Handler) TrackVisit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VisitRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, h.logger, "track", "TrackVisit", apperror.BadRequest("Invalid JSON"))
		return
	}
	if err := validate.Struct(req); err != nil {
		writeError(w, h.logger, "track", "TrackVisit", apperror.BadRequest("Validation failed"))
		return
	}

	ip := actionhttp.ExtractIP(r)
	action := domain.Action{
		VisitID:   req.VisitID,
		Source:    actionhttp.NormalizeSource(req.Source),
		Timestamp: req.Timestamp,
		IPAddress: ip,
	}

	_, err := h.trackService.TrackAction(r.Context(), "visit", action)
	if err != nil {
		writeError(w, h.logger, "track", "TrackVisit", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "visit recorded"})
}

// GetAllVisits возвращает все визиты.
// @Tags Visits
// @Summary Получить все визиты
// @Produce json
// @Param offset query int false "Смещение выборки (offset)"
// @Success 200 {array} domain.Action
// @Failure 404,500 {string} string
// @Router /track/visit/all [get]
func (h *Handler) GetAllVisits(w http.ResponseWriter, r *http.Request) {
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

	actions, err := h.visitService.GetAllVisits(r.Context(), h.limitAll, offset)
	if err != nil {
		writeError(w, h.logger, "track", "GetAllVisits", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(actions)
}

// GetStatsBySource возвращает статистику по источникам.
// @Tags Visits
// @Summary Статистика по источникам визитов
// @Produce json
// @Success 200 {array} visit.VisitSourceStat
// @Failure 404,500 {string} string
// @Router /track/stats/by-source [get]
func (h *Handler) GetStatsBySource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.visitService.GetStatsBySource(r.Context())
	if err != nil {
		writeError(w, h.logger, "track", "GetStatsBySource", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
