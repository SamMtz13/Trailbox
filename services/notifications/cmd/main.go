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

	notifctrl "trailbox/services/notifications/internal/controller"
	consulreg "trailbox/services/notifications/internal/discovery/consul"
	notifhttp "trailbox/services/notifications/internal/handler/http"
	notifmem "trailbox/services/notifications/internal/repository/memory"
)

const defaultPort = "8005"

func main() {
	// repo + controller
	repo := notifmem.NewRepository()
	ctrl := notifctrl.NewController(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", notifhttp.Healthz)

	h := notifhttp.NewNotificationHandler(ctrl)
	mux.HandleFunc("/v1/notifications", h.HandleList)      // GET
	mux.HandleFunc("/v1/notifications/send", h.HandleSend) // POST

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// consul register (estilo Registrar)
	registrar, err := consulreg.NewRegistrar()
	if err != nil {
		log.Printf("[notifications] consul init failed: %v", err)
	} else {
		svcName := getenv("SERVICE_NAME", "notifications")
		svcAddr := getenv("SERVICE_ADDRESS", "notifications")
		healthPath := getenv("SERVICE_HEALTH_PATH", "/health")
		_, err := registrar.Register(svcName, svcAddr, mustAtoi(port), healthPath)
		if err != nil {
			log.Printf("[notifications] consul register failed: %v", err)
		} else {
			log.Printf("[notifications] registered in consul as %s", svcName)
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
		log.Printf("[notifications] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[notifications] server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Printf("[notifications] shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[notifications] %s %s", r.Method, r.URL.Path)
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
