package http

import (
	"context"
	"encoding/json"
	"net/http"

	"trailbox/services/reviews/internal/controller"
)

type ReviewHandler struct {
	ctrl *controller.Controller
}

func NewReviewHandler(c *controller.Controller) *ReviewHandler {
	return &ReviewHandler{ctrl: c}
}

// GET /v1/reviews?route_id=&user_id=
func (h *ReviewHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	routeID := r.URL.Query().Get("route_id")
	userID := r.URL.Query().Get("user_id")
	list, err := h.ctrl.List(r.Context(), routeID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, list, http.StatusOK)
}

// POST /v1/reviews/create
// { "user_id":"...", "route_id":"...", "rating":5, "comment":"..." }
func (h *ReviewHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req controller.Review
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	res, err := h.ctrl.Create(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, res, http.StatusCreated)
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
