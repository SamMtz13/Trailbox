package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"trailbox/users/internal/controller/users"
	userhttp "trailbox/users/internal/handler/http"
	"trailbox/users/internal/repository/memory"
)

const defaultPort = "8001"

func main() {
	repo := memory.New()
	ctrl := users.NewController(repo)
	userHandler := userhttp.NewUserHandler(ctrl)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("users service ok"))
	})
	mux.HandleFunc("/users", userHandler.HandleList)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("[users] listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
