package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trailbox/services/gateway/internal/handler"
	"trailbox/services/gateway/internal/proxy"
)

const defaultPort = "8080"

func main() {
	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("/health", handler.Healthz)

	// Proxy hacia microservicios registrados
	mux.Handle("/users/", proxy.NewReverseProxy("http://users:50051"))
	mux.Handle("/routes/", proxy.NewReverseProxy("http://routes:50051"))
	mux.Handle("/workouts/", proxy.NewReverseProxy("http://workouts:50051"))
	mux.Handle("/reviews/", proxy.NewReverseProxy("http://reviews:50051"))
	mux.Handle("/leaderboard/", proxy.NewReverseProxy("http://leaderboard:50051"))
	mux.Handle("/notifications/", proxy.NewReverseProxy("http://notifications:50051"))

	port := getenvOr("PORT", defaultPort)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 🚀 Run Gateway
	go func() {
		log.Printf("[gateway] running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[gateway] server error: %v", err)
		}
	}()

	// 🛑 Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[gateway] shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("[gateway] shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[gateway] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
