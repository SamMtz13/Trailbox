package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "trailbox/gen/reviews"
	reviewsctrl "trailbox/services/reviews/internal/controller"
	reviewsdb "trailbox/services/reviews/internal/db"
	reviewsgrpc "trailbox/services/reviews/internal/handler/grpc"
	reviewrepo "trailbox/services/reviews/internal/repository/db"
)

const defaultPort = "50051"

func main() {
	// 1️⃣ DB + migración
	conn, err := reviewsdb.Connect()
	if err != nil {
		log.Fatalf("[reviews] ❌ DB error: %v", err)
	}
	log.Println("[reviews] ✅ Migración completada")

	repo := reviewrepo.New(conn)
	ctrl := reviewsctrl.NewController(repo)

	// 2️⃣ Servidor gRPC
	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[reviews] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterReviewsServer(grpcServer, reviewsgrpc.New(ctrl))

	// Health gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Servidor principal
	go func() {
		log.Printf("[reviews] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[reviews] server error: %v", err)
		}
	}()

	// 4️⃣ Apagado elegante
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[reviews] shutting down...")
	grpcServer.GracefulStop()

	sqlDB, _ := conn.DB()
	_ = sqlDB.Close()

	log.Println("[reviews] graceful shutdown complete")
}

// Helpers
func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
