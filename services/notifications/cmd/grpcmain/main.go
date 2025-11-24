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

	pb "trailbox/gen/notifications"
	notifctrl "trailbox/services/notifications/internal/controller"
	notifdb "trailbox/services/notifications/internal/db"
	notificationgrpc "trailbox/services/notifications/internal/handler/grpc"
	notifrepo "trailbox/services/notifications/internal/repository/db"
)

const defaultPort = "50051"

// ======================
// main()
// ======================
func main() {
	// 1️⃣ Conexión a DB
	conn, err := notifdb.Connect()
	if err != nil {
		log.Fatalf("[notifications] ❌ DB error: %v", err)
	}

	// Migración
	log.Println("[notifications] ✅ Migración completada")

	repo := notifrepo.New(conn)
	ctrl := notifctrl.NewController(repo)

	// 2️⃣ Servidor gRPC
	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[notifications] failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNotificationsServer(s, notificationgrpc.New(ctrl))

	// HealthCheck estándar gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(s, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Run server
	go func() {
		log.Printf("[notifications] 🚀 listening on :%s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[notifications] server error: %v", err)
		}
	}()

	// 4️⃣ Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[notifications] shutting down...")
	s.GracefulStop()

	sqlDB, _ := conn.DB()
	_ = sqlDB.Close()

	log.Println("[notifications] graceful shutdown complete")
}

// Helpers
func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
