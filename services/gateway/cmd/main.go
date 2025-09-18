package main

import (
	"log"
	"net/http"
	"os"

	"trailbox/services/gateway/internal/handler"
	"trailbox/services/gateway/internal/proxy"
)

const defaultPort = "8080"

func main() {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", handler.Healthz)

	// Proxy rules
	mux.Handle("/users/", proxy.NewReverseProxy("http://users:8001"))
	mux.Handle("/routes/", proxy.NewReverseProxy("http://routes:8002"))
	mux.Handle("/workouts/", proxy.NewReverseProxy("http://workouts:8003"))

	// Port config
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("gateway running on port %s", port)
	if err := http.ListenAndServe(":"+port, loggingMiddleware(mux)); err != nil {
		log.Fatalf("gateway error: %v", err)
	}
}

// Simple logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
