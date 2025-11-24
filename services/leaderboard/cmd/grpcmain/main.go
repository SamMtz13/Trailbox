package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "trailbox/gen/leaderboard"

	lbctrl "trailbox/services/leaderboard/internal/controller"
	lbdb "trailbox/services/leaderboard/internal/db"
	lbgrpc "trailbox/services/leaderboard/internal/handler/grpc"
	lbrepo "trailbox/services/leaderboard/internal/repository/db"
)

const defaultPort = "50051"

func main() {
	conn, err := lbdb.Connect()
	if err != nil {
		log.Fatalf("[leaderboard] DB error: %v", err)
	}

	repo := lbrepo.New(conn)
	ctrl := lbctrl.NewController(repo)

	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[leaderboard] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLeaderboardServer(grpcServer, lbgrpc.New(ctrl))

	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		log.Printf("[leaderboard] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[leaderboard] serve error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[leaderboard] shutting down...")
	grpcServer.GracefulStop()
	log.Println("[leaderboard] graceful shutdown complete")
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
