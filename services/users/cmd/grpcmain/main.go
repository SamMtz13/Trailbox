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

	pb "trailbox/gen/users"

	userctrl "trailbox/services/users/internal/controller/users"
	"trailbox/services/users/internal/db"
	usergrpc "trailbox/services/users/internal/handler/grpc"
	userrepo "trailbox/services/users/internal/repository/db"
)

const defaultPort = "50051"

func main() {
	// ===============================
	// 1️⃣ Conexión a PostgreSQL
	// ===============================
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("[users] ❌ error conectando a Postgres: %v", err)
	}

	// Extensión y migración
	log.Println("[users] ✅ Migración completada")

	// ===============================
	// 2️⃣ Controlador y repositorio
	// ===============================
	repo := userrepo.New(dbConn)
	ctrl := userctrl.NewController(repo)

	// ===============================
	// 3️⃣ gRPC setup
	// ===============================
	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[users] failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUsersServer(grpcServer, usergrpc.New(ctrl))

	// Health gRPC interno
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// ===============================
	// 4️⃣ Servidor gRPC principal
	// ===============================
	go func() {
		log.Printf("[users] 🚀 gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[users] server error: %v", err)
		}
	}()

	// ===============================
	// 7️⃣ Graceful shutdown
	// ===============================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[users] shutting down...")
	grpcServer.GracefulStop()
	log.Println("[users] graceful shutdown complete")
}

// Helpers
func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
