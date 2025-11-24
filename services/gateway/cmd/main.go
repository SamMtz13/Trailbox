package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	leaderboardpb "trailbox/gen/leaderboard"
	mappb "trailbox/gen/maps"
	notificationspb "trailbox/gen/notifications"
	reviewspb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	userspb "trailbox/gen/users"
	workoutspb "trailbox/gen/workouts"
	"trailbox/services/gateway/internal/handler"
)

const (
	defaultPort   = "8080"
	dialTimeout   = 10 * time.Second
	shutdownGrace = 10 * time.Second
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	usersConn := mustDialGRPC(getenvOr("USERS_SERVICE_ADDR", "users-service.final-project.svc.cluster.local:50051"))
	defer usersConn.Close()
	routesConn := mustDialGRPC(getenvOr("ROUTES_SERVICE_ADDR", "routes-service.final-project.svc.cluster.local:50051"))
	defer routesConn.Close()
	workoutsConn := mustDialGRPC(getenvOr("WORKOUTS_SERVICE_ADDR", "workouts-service.final-project.svc.cluster.local:50051"))
	defer workoutsConn.Close()
	reviewsConn := mustDialGRPC(getenvOr("REVIEWS_SERVICE_ADDR", "reviews-service.final-project.svc.cluster.local:50051"))
	defer reviewsConn.Close()
	leaderboardConn := mustDialGRPC(getenvOr("LEADERBOARD_SERVICE_ADDR", "leaderboard-service.final-project.svc.cluster.local:50051"))
	defer leaderboardConn.Close()
	notificationsConn := mustDialGRPC(getenvOr("NOTIFICATIONS_SERVICE_ADDR", "notifications-service.final-project.svc.cluster.local:50051"))
	defer notificationsConn.Close()
	mapsConn := mustDialGRPC(getenvOr("MAP_SERVICE_ADDR", "map-service.final-project.svc.cluster.local:50051"))
	defer mapsConn.Close()

	api := handler.NewAPI(
		userspb.NewUsersClient(usersConn),
		routespb.NewRoutesClient(routesConn),
		workoutspb.NewWorkoutsClient(workoutsConn),
		reviewspb.NewReviewsClient(reviewsConn),
		leaderboardpb.NewLeaderboardClient(leaderboardConn),
		notificationspb.NewNotificationsClient(notificationsConn),
		mappb.NewMapClient(mapsConn),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Healthz)
	api.Register(mux)

	port := getenvOr("PORT", defaultPort)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[gateway] running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[gateway] server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[gateway] shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[gateway] graceful shutdown error: %v", err)
	}
	log.Println("[gateway] shutdown complete")
}

func mustDialGRPC(target string) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("[gateway] failed to connect to %s: %v", target, err)
	}
	return conn
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[gateway] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := getenvOr("CORS_ALLOWED_ORIGINS", "*")
	allowedMethods := getenvOr("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	allowedHeaders := getenvOr("CORS_ALLOWED_HEADERS", "Content-Type,Authorization")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getenvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
