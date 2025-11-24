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

	pb "trailbox/gen/maps"

	mapctrl "trailbox/services/map/internal/controller"
	mapdb "trailbox/services/map/internal/db"
	mapgrpc "trailbox/services/map/internal/handler/grpc"
	maprepo "trailbox/services/map/internal/repository/db"
)

const defaultPort = "50051"

func main() {
	conn, err := mapdb.Connect()
	if err != nil {
		log.Fatalf("[map] DB error: %v", err)
	}

	repo := maprepo.New(conn)
	ctrl := mapctrl.NewController(repo)

	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[map] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMapServer(grpcServer, mapgrpc.New(ctrl))

	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		log.Printf("[map] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[map] serve error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[map] shutting down...")
	grpcServer.GracefulStop()
	log.Println("[map] graceful shutdown complete")
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
