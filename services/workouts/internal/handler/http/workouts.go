package http

import (
	"encoding/json"
	"net/http"
	"trailbox/services/workouts/internal/controller/workouts"
)

type WorkoutHandler struct {
	ctrl *workouts.Controller
}

func NewWorkoutHandler(c *workouts.Controller) *WorkoutHandler {
	return &WorkoutHandler{ctrl: c}
}
func (h *WorkoutHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.ctrl.ListWorkouts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
