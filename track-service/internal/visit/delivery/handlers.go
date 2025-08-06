package visithttp

import (
	"encoding/json"
	"net/http"

	"github.com/Vovarama1992/retry/pkg/domain"
	track "github.com/Vovarama1992/retry/track-service/internal/ports"
	visit_ports "github.com/Vovarama1992/retry/track-service/internal/visit/ports"
	validator "github.com/go-playground/validator/v10"
)

type Handler struct {
	trackService track.Service
	visitService visit_ports.VisitService
}

func NewHandler(trackService track.Service, visitService visit_ports.VisitService) *Handler {
	return &Handler{
		trackService: trackService,
		visitService: visitService,
	}
}

var validate = validator.New()

func (h *Handler) TrackVisit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VisitRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	ip := extractIP(r)

	action := domain.Action{
		VisitID:   req.VisitID,
		Source:    req.Source,
		Timestamp: req.Timestamp,
		IPAddress: ip,
	}

	_, err := h.trackService.TrackAction(r.Context(), "visit", action)
	if err != nil {
		http.Error(w, "Failed to track visit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "visit recorded"})
}

func (h *Handler) GetAllVisits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	actions, err := h.visitService.GetAllVisits(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch visits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actions)
}

func (h *Handler) GetStatsBySource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.visitService.GetStatsBySource(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
