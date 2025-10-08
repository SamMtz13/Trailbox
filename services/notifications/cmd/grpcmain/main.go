package main

import (
	"context"
	"log"
	"net"
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
	consulreg "trailbox/services/notifications/internal/discovery/consul"
	notifdb "trailbox/services/notifications/internal/repository/db"
)

const defaultPort = "50051"

// ======================
// gRPC Server Definition
// ======================
type notificationServer struct {
	pb.UnimplementedNotificationsServer
	ctrl *notifctrl.Controller
}

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
	port := getenvOr("PORT", defaultPort)

	// 1️⃣ Conexión a DB
	conn, err := notifdb.Connect()
	if err != nil {
		log.Fatalf("[notifications] DB error: %v", err)
	}

	repo := notifdb.New(conn)
	ctrl := notifctrl.NewController(repo)

	// 2️⃣ Servidor gRPC
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[notifications] failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNotificationsServer(s, &notificationServer{ctrl: ctrl})

	// HealthCheck estándar
	hs := health.NewServer()
	healthpb.RegisterHealthServer(s, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Registro en Consul
	reg, err := consulreg.NewRegistrar()
	if err != nil {
		log.Fatalf("[notifications] consul registrar init error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "notifications")
	healthPath := getenvOr("SERVICE_HEALTH_PATH", "/grpc.health.v1.Health/Check")

	id, err := reg.Register(getenvOr("SERVICE_NAME", "notifications"), addr, mustAtoi(port), healthPath)
	if err != nil {
		log.Fatalf("[notifications] consul register error: %v", err)
	}
	log.Printf("[notifications] consul registered id=%s", id)

	// 4️⃣ Run server
	go func() {
		log.Printf("[notifications] listening on :%s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[notifications] server error: %v", err)
		}
	}()

	// 5️⃣ Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("[notifications] shutting down...")
	s.GracefulStop()
	reg.Deregister()
	log.Println("[notifications] shutdown complete")
}

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
