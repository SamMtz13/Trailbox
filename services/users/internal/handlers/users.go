package handlers

import (
	"encoding/json"
	"net/http"

	"trailbox/services/users/internal/db"
	"trailbox/services/users/internal/model"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()
	var users []model.User

	if err := database.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
