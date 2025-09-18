package http

import (
	"encoding/json"
	"net/http"

	"trailbox/services/notifications/internal/controller"
)

type NotificationHandler struct{ ctrl *controller.Controller }

func NewNotificationHandler(c *controller.Controller) *NotificationHandler {
	return &NotificationHandler{ctrl: c}
}

// GET /v1/notifications
func (h *NotificationHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := h.ctrl.List(r.Context())
	if err != nil {
		http.Error(w, "cannot list notifications", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items, http.StatusOK)
}

// POST /v1/notifications/send
func (h *NotificationHandler) HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req controller.Notification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	res, err := h.ctrl.Send(r.Context(), req)
	if err != nil {
		http.Error(w, "cannot send", http.StatusInternalServerError)
		return
	}
	writeJSON(w, res, http.StatusCreated)
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
