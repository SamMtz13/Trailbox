package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"trailbox/services/users/internal/db"
	"trailbox/services/users/internal/handlers"
)

func main() {
	db.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Trailbox Users API - use /users to list users"}`))
	})
	http.HandleFunc("/users/", handlers.GetUsers)
	log.Printf("Users service listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
