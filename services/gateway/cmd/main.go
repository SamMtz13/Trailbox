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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"trailbox/services/gateway/internal/api"
	"trailbox/services/gateway/internal/handler"
	lbpb "trailbox/gen/leaderboard"
	mapspb "trailbox/gen/maps"
	notifpb "trailbox/gen/notifications"
	reviewpb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	userpb "trailbox/gen/users"
	workoutpb "trailbox/gen/workouts"
)

const defaultPort = "8080"

func main() {
	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("/health", handler.Healthz)

	conns := make([]*grpc.ClientConn, 0)
	dial := func(envKey, fallback string) *grpc.ClientConn {
		addr := getenvOr(envKey, fallback)
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("[gateway] failed to dial %s (%s): %v", envKey, addr, err)
		}
		conns = append(conns, conn)
		return conn
	}

	clients := api.Clients{
		Users:         userpb.NewUsersClient(dial("USERS_SERVICE_ADDR", defaultSvcAddr("users"))),
		Routes:        routespb.NewRoutesClient(dial("ROUTES_SERVICE_ADDR", defaultSvcAddr("routes"))),
		Workouts:      workoutpb.NewWorkoutsClient(dial("WORKOUTS_SERVICE_ADDR", defaultSvcAddr("workouts"))),
		Reviews:       reviewpb.NewReviewsClient(dial("REVIEWS_SERVICE_ADDR", defaultSvcAddr("reviews"))),
		Leaderboard:   lbpb.NewLeaderboardClient(dial("LEADERBOARD_SERVICE_ADDR", defaultSvcAddr("leaderboard"))),
		Notifications: notifpb.NewNotificationsClient(dial("NOTIFICATIONS_SERVICE_ADDR", defaultSvcAddr("notifications"))),
		Maps:          mapspb.NewMapClient(dial("MAPS_SERVICE_ADDR", defaultSvcAddr("maps"))),
	}

	api.New(clients).RegisterRoutes(mux)

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
	for _, conn := range conns {
		_ = conn.Close()
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
	return name + ".final-project.svc.cluster.local:50051"
}
