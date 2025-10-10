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
	http.HandleFunc("/users", handlers.GetUsers)

	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
