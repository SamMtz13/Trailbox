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

	mapctrl "trailbox/services/map/internal/controller"
	consulreg "trailbox/services/map/internal/discovery/consul"
	maphttp "trailbox/services/map/internal/handler/http"
	mapmem "trailbox/services/map/internal/repository/memory"
)

const defaultPort = "8006"

func main() {
	repo := mapmem.NewRepository()
	ctrl := mapctrl.NewController(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", maphttp.Healthz)

	h := maphttp.NewMapHandler(ctrl)
	mux.HandleFunc("/v1/maps/route", h.HandleGetRouteCoords)
	mux.HandleFunc("/v1/maps/route/set", h.HandleSetRouteCoords)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// consul register
	registrar, err := consulreg.NewRegistrar()
	if err != nil {
		log.Printf("[map] consul init failed: %v", err)
	} else {
		svcName := getenv("SERVICE_NAME", "map")
		svcAddr := getenv("SERVICE_ADDRESS", "map")
		healthPath := getenv("SERVICE_HEALTH_PATH", "/health")
		_, err := registrar.Register(svcName, svcAddr, mustAtoi(port), healthPath)
		if err != nil {
			log.Printf("[map] consul register failed: %v", err)
		} else {
			log.Printf("[map] registered in consul as %s", svcName)
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
		log.Printf("[map] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[map] server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Printf("[map] shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[map] %s %s", r.Method, r.URL.Path)
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
