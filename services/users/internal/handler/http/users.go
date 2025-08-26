package http

import (
	"encoding/json"
	"net/http"
	"trailbox/users/internal/controller/users"
)

type UserHandler struct {
	ctrl *users.Controller
}

func NewUserHandler(c *users.Controller) *UserHandler {
	return &UserHandler{ctrl: c}
}

func (h *UserHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.ctrl.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
