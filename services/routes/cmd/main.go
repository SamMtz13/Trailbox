package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"trailbox/routes/internal/controller/routes"
	routehttp "trailbox/routes/internal/handler/http"
	"trailbox/routes/internal/repository/memory"
)

const defaultPort = "8002"

func main() {
	repo := memory.New()
	ctrl := routes.NewController(repo)
	routesHandler := routehttp.NewRouteHandler(ctrl)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("routes service ok"))
	})
	mux.HandleFunc("/routes", routesHandler.HandleList)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("[routes] listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
