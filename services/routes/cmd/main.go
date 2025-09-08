package main

import (
	"log"
	"net/http"
	"os"
	"time"

	routesctrl "trailbox/routes/internal/controller/routes"
	routehttp "trailbox/routes/internal/handler/http"
	routememory "trailbox/routes/internal/repository/memory"
)

const defaultPort = "8002"

func main() {
	// Repository + controller
	repo := routememory.NewRepository()
	ctrl := routesctrl.NewController(repo)

	// HTTP mux
	mux := http.NewServeMux()

	// Handlers existentes
	mux.HandleFunc("/health", routehttp.Healthz)
	routeHandler := routehttp.NewRouteHandler(ctrl)
	mux.HandleFunc("/v1/routes", routeHandler.HandleList)

	// Handlers nuevos para hablar con el peer (curl en Go)
	peerHandler := routehttp.NewPeerHandler()
	mux.HandleFunc("/peer/health", peerHandler.HandlePeerHealth) // llama {PEER_URL}/health
	mux.HandleFunc("/peer/proxy", peerHandler.HandlePeerProxy)   // llama {PEER_URL}?path=/...

	// Heartbeat opcional cada 30s (desactiva con DISABLE_HEARTBEAT=1)
	if os.Getenv("DISABLE_HEARTBEAT") == "" {
		go func() {
			stop := make(chan struct{})
			defer close(stop)
			hb := routehttp.NewHeartbeat(peerHandler, "/health", 30*time.Second)
			hb.Run(stop, func(format string, args ...any) { log.Printf(format, args...) })
		}()
	}

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
