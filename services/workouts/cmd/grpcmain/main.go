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

	pb "trailbox/gen/workouts"

	wctrl "trailbox/services/workouts/internal/controller/workouts"
	"trailbox/services/workouts/internal/db"
	workoutgrpc "trailbox/services/workouts/internal/handler/grpc"
	wrepo "trailbox/services/workouts/internal/repository/db"
)

const defaultPort = "50051"

func main() {
	// 1️⃣ Conexión DB + migración
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("[workouts] ❌ DB connection error: %v", err)
	}
	log.Println("[workouts] ✅ Migración completada")

	repo := wrepo.New(conn)
	ctrl := wctrl.NewController(repo)

	// 2️⃣ Servidor gRPC
	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[workouts] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWorkoutsServer(grpcServer, workoutgrpc.New(ctrl))

	// Health gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Arranque del servidor
	go func() {
		log.Printf("[workouts] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[workouts] server error: %v", err)
		}
	}()

	// 4️⃣ Apagado elegante
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[workouts] shutting down...")
	grpcServer.GracefulStop()

	sqlDB, _ := conn.DB()
	_ = sqlDB.Close()

	log.Println("[workouts] graceful shutdown complete")
}

// Helpers
func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
