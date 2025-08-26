package http

import (
	"encoding/json"
	"net/http"
	"trailbox/routes/internal/controller/routes"
)

type RouteHandler struct {
	ctrl *routes.Controller
}

func NewRouteHandler(c *routes.Controller) *RouteHandler {
	return &RouteHandler{ctrl: c}
}

func (h *RouteHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.ctrl.ListRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
