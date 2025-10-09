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
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	commonpb "trailbox/gen/common"
	pb "trailbox/gen/workouts"

	wctrl "trailbox/services/workouts/internal/controller/workouts"
	"trailbox/services/workouts/internal/db"
	wconsul "trailbox/services/workouts/internal/discovery/consul"
	"trailbox/services/workouts/internal/model"
	wrepo "trailbox/services/workouts/internal/repository/db"

	"github.com/joho/godotenv"
)

const (
	defaultPort    = "50051"
	healthHTTPPort = 8081
)

type workoutServer struct {
	pb.UnimplementedWorkoutsServer
	ctrl *wctrl.Controller
}

// GetWorkout obtiene un workout por ID
func (s *workoutServer) GetWorkout(ctx context.Context, req *commonpb.UserId) (*pb.Workout, error) {
	w, err := s.ctrl.GetWorkout(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "workout not found")
	}
	return &pb.Workout{
		Id:       w.ID.String(),
		UserId:   w.UserID.String(),
		RouteId:  w.RouteID.String(),
		Date:     w.Date.Format(time.RFC3339),
		Duration: float64(w.Duration),
		Calories: float64(w.Calories),
	}, nil
}

// ListWorkouts lista todos los workouts
func (s *workoutServer) ListWorkouts(ctx context.Context, req *pb.ListWorkoutsRequest) (*pb.ListWorkoutsResponse, error) {
	workouts, err := s.ctrl.ListWorkouts()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list workouts")
	}
	resp := &pb.ListWorkoutsResponse{}
	for _, w := range workouts {
		resp.Workouts = append(resp.Workouts, &pb.Workout{
			Id:       w.ID.String(),
			UserId:   w.UserID.String(),
			RouteId:  w.RouteID.String(),
			Date:     w.Date.Format(time.RFC3339),
			Duration: float64(w.Duration),
			Calories: float64(w.Calories),
		})
	}
	return resp, nil
}

func main() {
	_ = godotenv.Load()

	// 1️⃣ Conexión DB + migración
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("[workouts] ❌ DB connection error: %v", err)
	}

	if err := conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("[workouts] ❌ failed to create uuid extension: %v", err)
	}

	if err := conn.AutoMigrate(&model.Workout{}); err != nil {
		log.Fatalf("[workouts] ❌ migration error: %v", err)
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
	pb.RegisterWorkoutsServer(grpcServer, &workoutServer{ctrl: ctrl})

	// Health gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Health HTTP adicional (para Consul)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK - workouts"))
		})
		log.Printf("[workouts] health HTTP listening on :%d", healthHTTPPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", healthHTTPPort), nil); err != nil {
			log.Printf("[workouts] health server error: %v", err)
		}
	}()

	// 4️⃣ Registro en Consul
	reg, err := wconsul.NewRegistrar()
	if err != nil {
		log.Fatalf("[workouts] consul init error: %v", err)
	}

	addr := getenvOr("SERVICE_ADDRESS", "workouts")
	id, err := reg.Register(getenvOr("SERVICE_NAME", "workouts"), addr, healthHTTPPort, "/health")
	if err != nil {
		log.Fatalf("[workouts] consul register error: %v", err)
	}
	log.Printf("[workouts] consul registered id=%s", id)

	// 5️⃣ Arranque del servidor
	go func() {
		log.Printf("[workouts] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[workouts] server error: %v", err)
		}
	}()

	// 6️⃣ Apagado elegante
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[workouts] shutting down...")
	grpcServer.GracefulStop()
	reg.Deregister()

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

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
