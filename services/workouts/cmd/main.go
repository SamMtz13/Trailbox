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

	wctrl "trailbox/services/workouts/internal/controller/workouts"
	wconsul "trailbox/services/workouts/internal/discovery/consul"
	wh "trailbox/services/workouts/internal/handler/http"
	wmem "trailbox/services/workouts/internal/repository/memory"
)

const defaultPort = "8003"

func main() {
	// Core wiring
	repo := wmem.New()
	ctrl := wctrl.NewController(repo)
	mux := http.NewServeMux()

	// Health + endpoints
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("workouts service ok"))
	})
	whandler := wh.NewWorkoutHandler(ctrl)
	mux.HandleFunc("/workouts", whandler.HandleList)

	// Port
	port := getenvOr("PORT", defaultPort)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ===== Consul register BEFORE serving =====
	reg, err := wconsul.NewRegistrar()
	if err != nil {
		log.Fatalf("[workouts] consul registrar init error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "workouts")
	healthPath := getenvOr("SERVICE_HEALTH_PATH", "/healthz")
	id, err := reg.Register(getenvOr("SERVICE_NAME", "workouts"), addr, mustAtoi(port), healthPath)
	if err != nil {
		log.Fatalf("[workouts] consul register error: %v", err)
	}
	log.Printf("[workouts] consul registered id=%s", id)
	// =========================================

	// Serve
	go func() {
		log.Printf("[workouts] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[workouts] server error: %v", err)
		}
	}()

	// Graceful shutdown + deregister
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	reg.Deregister()
	log.Println("[workouts] graceful shutdown complete")
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
