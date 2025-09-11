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

	routesctrl "trailbox/services/routes/internal/controller/routes"
	consuldiscovery "trailbox/services/routes/internal/discovery/consul"
	routehttp "trailbox/services/routes/internal/handler/http"
	routememory "trailbox/services/routes/internal/repository/memory"
)

const defaultPort = "8002"

func main() {
	// Repository + controller
	repo := routememory.NewRepository()
	ctrl := routesctrl.NewController(repo)

	// HTTP mux
	mux := http.NewServeMux()
	mux.HandleFunc("/health", routehttp.Healthz) // asegúrate que SERVICE_HEALTH_PATH=/health
	routeHandler := routehttp.NewRouteHandler(ctrl)
	mux.HandleFunc("/v1/routes", routeHandler.HandleList)

	// Peer tools
	peerHandler := routehttp.NewPeerHandler()
	mux.HandleFunc("/peer/health", peerHandler.HandlePeerHealth)
	// mux.HandleFunc("/peer/proxy", peerHandler.HandlePeerProxy)

	// Heartbeat opcional cada 30s
	if os.Getenv("DISABLE_HEARTBEAT") == "" {
		go func() {
			stop := make(chan struct{})
			defer close(stop)
			hb := routehttp.NewHeartbeat(peerHandler, "/health", 30*time.Second)
			hb.Run(stop, func(format string, args ...any) { log.Printf(format, args...) })
		}()
	}

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ===== Registro en Consul (ANTES de levantar el server) =====
	reg, err := consuldiscovery.NewRegistrar()
	if err != nil {
		log.Fatalf("[consul] registrar init error: %v", err)
	}

	addr := getenvOr("SERVICE_ADDRESS", "routes") // hostname del contenedor
	healthPath := getenvOr("SERVICE_HEALTH_PATH", "/health")
	portNum := mustAtoi(port)

	id, err := reg.Register(getenvOr("SERVICE_NAME", "routes"), addr, portNum, healthPath)
	if err != nil {
		log.Fatalf("[consul] register error: %v", err)
	}
	log.Printf("[consul] service registered id=%s", id)
	// ============================================================

	// Levantar servidor en goroutine
	go func() {
		log.Printf("[routes] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[routes] server error: %v", err)
		}
	}()

	// Esperar señal para apagar
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Shutdown gracioso + desregistro
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	reg.Deregister()
	log.Println("[routes] graceful shutdown complete")
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
