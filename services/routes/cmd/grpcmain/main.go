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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	commonpb "trailbox/gen/common"
	pb "trailbox/gen/routes"

	routesctrl "trailbox/services/routes/internal/controller/routes"
	"trailbox/services/routes/internal/db"
	"trailbox/services/routes/internal/model"
	routesrepo "trailbox/services/routes/internal/repository/db"

	"github.com/joho/godotenv"
)

type routeServer struct {
	pb.UnimplementedRoutesServer
	ctrl *routesctrl.Controller
}

// Obtener una ruta por ID
func (s *routeServer) GetRoute(ctx context.Context, req *commonpb.RouteId) (*pb.Route, error) {
	route, err := s.ctrl.GetRoute(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "route not found")
	}
	return &pb.Route{
		Id:            route.ID.String(),
		Name:          route.Path,
		DistanceKm:    float64(route.Distance),
		ElevationGain: float64(route.Duration),
	}, nil
}

// Listar todas las rutas
func (s *routeServer) ListRoutes(ctx context.Context, req *pb.ListRoutesRequest) (*pb.ListRoutesResponse, error) {
	routes, err := s.ctrl.ListRoutes()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list routes")
	}
	resp := &pb.ListRoutesResponse{}
	for _, r := range routes {
		resp.Routes = append(resp.Routes, &pb.Route{
			Id:            r.ID.String(),
			Name:          r.Path,
			DistanceKm:    float64(r.Distance),
			ElevationGain: float64(r.Duration),
		})
	}
	return resp, nil
}

const (
	defaultPort    = "50051"
	healthHTTPPort = 8081
)

func main() {
	_ = godotenv.Load()

	// ===============================
	// 1️⃣ Conexión DB + migración
	// ===============================
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("[routes] ❌ DB connection error: %v", err)
	}

	if err := conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("[routes] ❌ failed to create uuid extension: %v", err)
	}

	if err := conn.AutoMigrate(&model.Route{}); err != nil {
		log.Fatalf("[routes] ❌ migration error: %v", err)
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
	pb.RegisterRoutesServer(grpcServer, &routeServer{ctrl: ctrl})

	// Health gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// ===============================
	// 4️⃣ HTTP health (for k8s probes)
	// ===============================
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK - routes"))
		})
		log.Printf("[routes] health HTTP server on :%d", healthHTTPPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", healthHTTPPort), nil); err != nil {
			log.Fatalf("[routes] health server error: %v", err)
		}
	}()

	log.Printf("[routes] readiness HTTP on :%d", healthHTTPPort)

	// ===============================
	// 5️⃣ Arranque del servidor
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

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
