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

	reviewsctrl "trailbox/services/reviews/internal/controller"
	consulreg "trailbox/services/reviews/internal/discovery/consul"
	reviewhttp "trailbox/services/reviews/internal/handler/http"
	reviewmem "trailbox/services/reviews/internal/repository/memory"
)

const defaultPort = "8004"

func main() {
	// repo + controller
	repo := reviewmem.NewRepository()
	ctrl := reviewsctrl.NewController(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", reviewhttp.Healthz)

	h := reviewhttp.NewReviewHandler(ctrl)
	mux.HandleFunc("/v1/reviews", h.HandleList)          // GET ?route_id=&user_id=
	mux.HandleFunc("/v1/reviews/create", h.HandleCreate) // POST JSON

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// consul register (best effort, estilo actual con Registrar)
	registrar, err := consulreg.NewRegistrar()
	if err != nil {
		log.Printf("[reviews] consul init failed: %v", err)
	} else {
		svcName := getenv("SERVICE_NAME", "reviews")
		svcAddr := getenv("SERVICE_ADDRESS", "reviews")
		healthPath := getenv("SERVICE_HEALTH_PATH", "/health")
		_, err := registrar.Register(svcName, svcAddr, mustAtoi(port), healthPath)
		if err != nil {
			log.Printf("[reviews] consul register failed: %v", err)
		} else {
			log.Printf("[reviews] registered in consul as %s", svcName)
		}
		defer registrar.Deregister() // 🔹 deregister on exit
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("[reviews] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[reviews] server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Printf("[reviews] shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[reviews] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
