package http

import (
	"encoding/json"
	"net/http"

	"trailbox/services/map/internal/controller"
)

type MapHandler struct{ ctrl *controller.Controller }

func NewMapHandler(c *controller.Controller) *MapHandler { return &MapHandler{ctrl: c} }

func (h *MapHandler) HandleGetRouteCoords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	rm, err := h.ctrl.GetRoute(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, rm, http.StatusOK)
}

func (h *MapHandler) HandleSetRouteCoords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var rm controller.RouteMap
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.ctrl.SetRoute(r.Context(), rm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusCreated)
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
