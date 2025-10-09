package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "trailbox/gen/notifications"
	notifctrl "trailbox/services/notifications/internal/controller"
	notifdb "trailbox/services/notifications/internal/db"
	notifconsul "trailbox/services/notifications/internal/discovery/consul"
	"trailbox/services/notifications/internal/model"
	notifrepo "trailbox/services/notifications/internal/repository/db"

	"github.com/joho/godotenv"
)

const (
	defaultPort    = "50051"
	healthHTTPPort = 8081
)

// ======================
// gRPC Server Definition
// ======================
type notificationServer struct {
	pb.UnimplementedNotificationsServer
	ctrl *notifctrl.Controller
}

// 🔹 Obtener notificaciones por usuario
func (s *notificationServer) GetNotifications(ctx context.Context, req *pb.UserIdRequest) (*pb.NotificationsResponse, error) {
	notifs, err := s.ctrl.ListByUser(req.UserId)
	if err != nil {
		return nil, err
	}

	resp := &pb.NotificationsResponse{}
	for _, n := range notifs {
		resp.Notifications = append(resp.Notifications, &pb.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			Message:   n.Message,
			Read:      n.Read,
			CreatedAt: n.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

// 🔹 Enviar una nueva notificación
func (s *notificationServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.Notification, error) {
	n, err := s.ctrl.Create(req.UserId, req.Message)
	if err != nil {
		return nil, err
	}
	return &pb.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Message:   n.Message,
		Read:      n.Read,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ======================
// main()
// ======================
func main() {
	_ = godotenv.Load()

	// 1️⃣ Conexión a DB
	conn, err := notifdb.Connect()
	if err != nil {
		log.Fatalf("[notifications] ❌ DB error: %v", err)
	}

	// Migración
	if err := conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("[notifications] ❌ error creando extensión uuid-ossp: %v", err)
	}
	if err := conn.AutoMigrate(&model.Notification{}); err != nil {
		log.Fatalf("[notifications] ❌ error migrando modelo: %v", err)
	}
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
	pb.RegisterNotificationsServer(s, &notificationServer{ctrl: ctrl})

	// HealthCheck estándar gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(s, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Health HTTP adicional (para Consul)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK - notifications")
		})
		log.Printf("[notifications] health HTTP listening on :%d", healthHTTPPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", healthHTTPPort), nil); err != nil {
			log.Printf("[notifications] health HTTP error: %v", err)
		}
	}()

	// 4️⃣ Registro en Consul
	reg, err := notifconsul.NewRegistrar()
	if err != nil {
		log.Fatalf("[notifications] consul registrar init error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "notifications")

	id, err := reg.Register(getenvOr("SERVICE_NAME", "notifications"), addr, healthHTTPPort, "/health")
	if err != nil {
		log.Fatalf("[notifications] consul register error: %v", err)
	}
	log.Printf("[notifications] consul registered id=%s", id)

	// 5️⃣ Run server
	go func() {
		log.Printf("[notifications] 🚀 listening on :%s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[notifications] server error: %v", err)
		}
	}()

	// 6️⃣ Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[notifications] shutting down...")
	s.GracefulStop()
	reg.Deregister()

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

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
