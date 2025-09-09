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

	userctrl "trailbox/services/users/internal/controller/users"
	userconsul "trailbox/services/users/internal/discovery/consul"
	userhttp "trailbox/services/users/internal/handler/http"
	usermemory "trailbox/services/users/internal/repository/memory"
)

const defaultPort = "8001"

func main() {
	// Core app wiring
	repo := usermemory.New()
	ctrl := userctrl.NewController(repo)
	mux := http.NewServeMux()

	// Health + endpoints
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("users service ok"))
	})
	handler := userhttp.NewUserHandler(ctrl)
	mux.HandleFunc("/users", handler.HandleList)

	// Port
	port := getenvOr("PORT", defaultPort)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ===== Consul register BEFORE serving =====
	reg, err := userconsul.NewRegistrar()
	if err != nil {
		log.Fatalf("[users] consul registrar init error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "users")
	healthPath := getenvOr("SERVICE_HEALTH_PATH", "/healthz")
	id, err := reg.Register(getenvOr("SERVICE_NAME", "users"), addr, mustAtoi(port), healthPath)
	if err != nil {
		log.Fatalf("[users] consul register error: %v", err)
	}
	log.Printf("[users] consul registered id=%s", id)
	// =========================================

	// Serve in goroutine
	go func() {
		log.Printf("[users] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[users] server error: %v", err)
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
	log.Println("[users] graceful shutdown complete")
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
