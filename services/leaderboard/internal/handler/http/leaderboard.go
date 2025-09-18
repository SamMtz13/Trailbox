package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"trailbox/services/leaderboard/internal/controller"
)

type LeaderboardHandler struct{ ctrl *controller.Controller }

func NewLeaderboardHandler(c *controller.Controller) *LeaderboardHandler {
	return &LeaderboardHandler{ctrl: c}
}

func (h *LeaderboardHandler) HandleTop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	limit := 0
	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			limit = n
		}
	}
	items, err := h.ctrl.Top(r.Context(), limit)
	if err != nil {
		http.Error(w, "cannot get leaderboard", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items, http.StatusOK)
}

func (h *LeaderboardHandler) HandleUpsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		UserID     string  `json:"user_id"`
		DistanceKM float64 `json:"distance_km"`
		Workouts   int     `json:"workouts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	e, err := h.ctrl.Upsert(r.Context(), body.UserID, body.DistanceKM, body.Workouts)
	if err != nil {
		http.Error(w, "cannot upsert", http.StatusInternalServerError)
		return
	}
	writeJSON(w, e, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
