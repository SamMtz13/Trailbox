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

	pb "trailbox/gen/routes"

	routesctrl "trailbox/services/routes/internal/controller/routes"
	"trailbox/services/routes/internal/db"
	routesgrpc "trailbox/services/routes/internal/handler/grpc"
	routesrepo "trailbox/services/routes/internal/repository/db"
)

const defaultPort = "50051"

func main() {
	// ===============================
	// 1️⃣ Conexión DB + migración
	// ===============================
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("[routes] ❌ DB connection error: %v", err)
	}
	log.Println("[routes] ✅ Migración completada")

	// ===============================
	// 2️⃣ Repo + controller
	// ===============================
	repo := routesrepo.New(conn)
	ctrl := routesctrl.NewController(repo)

	// ===============================
	// 3️⃣ Servidor gRPC
	// ===============================
	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[routes] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRoutesServer(grpcServer, routesgrpc.New(ctrl))

	// Health gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// ===============================
	// 4️⃣ Arranque del servidor
	// ===============================
	go func() {
		log.Printf("[routes] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[routes] server error: %v", err)
		}
	}()

	// ===============================
	// 6️⃣ Apagado elegante
	// ===============================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[routes] shutting down...")
	grpcServer.GracefulStop()
	log.Println("[routes] graceful shutdown complete")
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
