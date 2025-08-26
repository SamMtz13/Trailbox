package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"trailbox/workouts/internal/controller/workouts"
	workouthttp "trailbox/workouts/internal/handler/http"
	"trailbox/workouts/internal/repository/memory"
)

const defaultPort = "8003"

func main() {
	repo := memory.New()
	ctrl := workouts.NewController(repo)
	workoutHandler := workouthttp.NewWorkoutHandler(ctrl)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("workouts service ok"))
	})
	mux.HandleFunc("/workouts", workoutHandler.HandleList)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("[workouts] listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}