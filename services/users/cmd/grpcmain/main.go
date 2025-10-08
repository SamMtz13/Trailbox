package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	commonpb "trailbox/gen/common"
	pb "trailbox/gen/users"

	userctrl "trailbox/services/users/internal/controller/users"
	"trailbox/services/users/internal/db"
	userconsul "trailbox/services/users/internal/discovery/consul"
	"trailbox/services/users/internal/model"
	userrepo "trailbox/services/users/internal/repository/db"

	"github.com/joho/godotenv"
)

type userServer struct {
	pb.UnimplementedUsersServer
	ctrl *userctrl.Controller
}

// GetUser usa el controller para responder
func (s *userServer) GetUser(ctx context.Context, req *commonpb.UserId) (*pb.User, error) {
	user, err := s.ctrl.GetUser(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &pb.User{
		Id:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// ListUsers retorna todos los usuarios
func (s *userServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.ctrl.ListUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}
	resp := &pb.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &pb.User{
			Id:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return resp, nil
}

const defaultPort = "50051"

func main() {
	// Cargar variables de entorno
	_ = godotenv.Load()

	// ===============================
	// 🔹 1. Conexión a PostgreSQL
	// ===============================
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("[users] ❌ error conectando a Postgres: %v", err)
	}

	// Crear extensión para UUID si no existe
	if err := dbConn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("[users] ❌ error creando extensión uuid-ossp: %v", err)
	}

	// Migrar modelo users
	if err := dbConn.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("[users] ❌ error migrando modelo users: %v", err)
	}
	log.Println("[users] ✅ Migración de users completada")

	// ===============================
	// 🔹 2. Inicializar repositorio y controlador
	// ===============================
	repo := userrepo.New(dbConn)
	ctrl := userctrl.NewController(repo)

	// ===============================
	// 🔹 3. gRPC setup
	// ===============================
	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[users] failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUsersServer(grpcServer, &userServer{ctrl: ctrl})

	// Health gRPC estándar
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// ===============================
	// 🔹 4. Registro en Consul
	// ===============================
	reg, err := userconsul.NewRegistrar()
	if err != nil {
		log.Fatalf("[users] consul registrar init error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "users")
	healthPath := getenvOr("SERVICE_HEALTH_PATH", "/grpc.health.v1.Health/Check")
	id, err := reg.Register(getenvOr("SERVICE_NAME", "users"), addr, mustAtoi(port), healthPath)
	if err != nil {
		log.Fatalf("[users] consul register error: %v", err)
	}
	log.Printf("[users] consul registered id=%s", id)

	// ===============================
	// 🔹 5. Servidor gRPC
	// ===============================
	go func() {
		log.Printf("[users] 🚀 gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[users] server error: %v", err)
		}
	}()

	// ===============================
	// 🔹 6. Graceful shutdown
	// ===============================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[users] shutting down...")
	grpcServer.GracefulStop()
	reg.Deregister()
	log.Println("[users] graceful shutdown complete")
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
