package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	aggcontroller "trailbox/services/gateway/internal/aggregator/controller"
	gatewayclients "trailbox/services/gateway/internal/clients"
	gatewayleaderboard "trailbox/services/gateway/internal/gateway/leaderboard/grpc"
	gatewaymaps "trailbox/services/gateway/internal/gateway/maps/grpc"
	gatewaynotifications "trailbox/services/gateway/internal/gateway/notifications/grpc"
	gatewayreviews "trailbox/services/gateway/internal/gateway/reviews/grpc"
	gatewayroutes "trailbox/services/gateway/internal/gateway/routes/grpc"
	gatewayusers "trailbox/services/gateway/internal/gateway/users/grpc"
	gatewayworkouts "trailbox/services/gateway/internal/gateway/workouts/grpc"
	gatewayhttp "trailbox/services/gateway/internal/http/handler"
)

const defaultPort = "8080"

func main() {
	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	closers := make([]interface{ Close() error }, 0)

	mustDialUsers := func(envKey, fallback string) *gatewayusers.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewayusers.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial users (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	mustDialRoutes := func(envKey, fallback string) *gatewayroutes.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewayroutes.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial routes (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	mustDialWorkouts := func(envKey, fallback string) *gatewayworkouts.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewayworkouts.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial workouts (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	mustDialReviews := func(envKey, fallback string) *gatewayreviews.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewayreviews.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial reviews (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	mustDialLeaderboard := func(envKey, fallback string) *gatewayleaderboard.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewayleaderboard.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial leaderboard (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	mustDialNotifications := func(envKey, fallback string) *gatewaynotifications.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewaynotifications.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial notifications (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	mustDialMaps := func(envKey, fallback string) *gatewaymaps.Client {
		addr := getenvOr(envKey, fallback)
		client, err := gatewaymaps.Dial(addr)
		if err != nil {
			log.Fatalf("[gateway] failed to dial maps (%s): %v", addr, err)
		}
		closers = append(closers, client)
		return client
	}

	usersClient := mustDialUsers("USERS_SERVICE_ADDR", defaultSvcAddr("users"))
	routesClient := mustDialRoutes("ROUTES_SERVICE_ADDR", defaultSvcAddr("routes"))
	workoutsClient := mustDialWorkouts("WORKOUTS_SERVICE_ADDR", defaultSvcAddr("workouts"))
	reviewsClient := mustDialReviews("REVIEWS_SERVICE_ADDR", defaultSvcAddr("reviews"))
	leaderboardClient := mustDialLeaderboard("LEADERBOARD_SERVICE_ADDR", defaultSvcAddr("leaderboard"))
	notificationsClient := mustDialNotifications("NOTIFICATIONS_SERVICE_ADDR", defaultSvcAddr("notifications"))
	mapsClient := mustDialMaps("MAPS_SERVICE_ADDR", defaultSvcAddr("maps"))

	clientSet := gatewayclients.Clients{
		Users:         usersClient.API(),
		Routes:        routesClient.API(),
		Workouts:      workoutsClient.API(),
		Reviews:       reviewsClient.API(),
		Leaderboard:   leaderboardClient.API(),
		Notifications: notificationsClient.API(),
		Maps:          mapsClient.API(),
	}

	aggregatorController := aggcontroller.New(clientSet)
	gatewayhttp.New(clientSet, aggregatorController).Register(mux)

	port := getenvOr("PORT", defaultPort)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(corsMiddleware(mux)),
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
	for _, c := range closers {
		_ = c.Close()
	}
	log.Println("[gateway] shutdown complete")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[gateway] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func defaultSvcAddr(name string) string {
	if strings.Contains(name, ".") {
		return name
	}
	return name + ".default.svc.cluster.local:50051"
}
