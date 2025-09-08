package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"trailbox/routes/internal/client/peer"
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

type userDTO struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

type listWithUsersResp struct {
	Routes any       `json:"routes"`
	Users  []userDTO `json:"users"`
}

// GET /v1/routes-with-users -> combina lista de rutas + users desde el microservicio users
func (h *RouteHandler) HandleListWithUsers(w http.ResponseWriter, r *http.Request) {
	// 1) Obtener rutas locales
	routesList, err := h.ctrl.ListRoutes()
	if err != nil {
		http.Error(w, "failed to list routes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) Llamar a users: GET {PEER_URL}/v1/users
	c, err := peer.New()
	if err != nil {
		http.Error(w, "peer not configured (set PEER_URL)", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	code, body, err := c.Get(ctx, "/v1/users")
	if err != nil {
		http.Error(w, "peer call failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	if code >= 400 {
		http.Error(w, "peer returned status "+http.StatusText(code), http.StatusBadGateway)
		return
	}

	var users []userDTO
	if err := json.Unmarshal(body, &users); err != nil {
		http.Error(w, "invalid users JSON: "+err.Error(), http.StatusBadGateway)
		return
	}

	// 3) Responder combinado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listWithUsersResp{
		Routes: routesList,
		Users:  users,
	})
}
