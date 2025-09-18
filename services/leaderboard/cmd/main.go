package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	lbctrl "trailbox/services/leaderboard/internal/controller"
	consulreg "trailbox/services/leaderboard/internal/discovery/consul"
	lbhttp "trailbox/services/leaderboard/internal/handler/http"
	lbmem "trailbox/services/leaderboard/internal/repository/memory"
)

const defaultPort = "8007"

func main() {
	repo := lbmem.NewRepository()
	ctrl := lbctrl.NewController(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", lbhttp.Healthz)

	h := lbhttp.NewLeaderboardHandler(ctrl)
	mux.HandleFunc("/v1/leaderboard", h.HandleTop)
	mux.HandleFunc("/v1/leaderboard/update", h.HandleUpsert)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// consul register
	registrar, err := consulreg.NewRegistrar()
	if err != nil {
		log.Printf("[leaderboard] consul init failed: %v", err)
	} else {
		svcName := getenv("SERVICE_NAME", "leaderboard")
		svcAddr := getenv("SERVICE_ADDRESS", "leaderboard")
		healthPath := getenv("SERVICE_HEALTH_PATH", "/health")
		_, err := registrar.Register(svcName, svcAddr, mustAtoi(port), healthPath)
		if err != nil {
			log.Printf("[leaderboard] consul register failed: %v", err)
		} else {
			log.Printf("[leaderboard] registered in consul as %s", svcName)
		}
		defer registrar.Deregister()
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("[leaderboard] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[leaderboard] server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Printf("[leaderboard] shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[leaderboard] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustAtoi(s string) int { n, _ := strconv.Atoi(s); return n }
