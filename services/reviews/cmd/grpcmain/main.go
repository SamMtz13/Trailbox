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

	pb "trailbox/gen/reviews"
	reviewsctrl "trailbox/services/reviews/internal/controller"
	reviewsdb "trailbox/services/reviews/internal/db"
	"trailbox/services/reviews/internal/model"
	reviewrepo "trailbox/services/reviews/internal/repository/db"

	"github.com/joho/godotenv"
)

const (
	defaultPort    = "50051"
	healthHTTPPort = 8081
)

type reviewServer struct {
	pb.UnimplementedReviewsServer
	ctrl *reviewsctrl.Controller
}

// 🔹 Obtener lista de reseñas
func (s *reviewServer) GetReviews(ctx context.Context, req *pb.ReviewListRequest) (*pb.ReviewListResponse, error) {
	revs, err := s.ctrl.ListReviews(req.RouteId)
	if err != nil {
		return nil, err
	}

	resp := &pb.ReviewListResponse{}
	for _, r := range revs {
		resp.Reviews = append(resp.Reviews, &pb.Review{
			Id:        r.ID,
			UserId:    r.UserID,
			RouteId:   r.RouteID,
			Rating:    int32(r.Rating),
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

// 🔹 Crear una reseña nueva
func (s *reviewServer) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.Review, error) {
	r, err := s.ctrl.AddReview(req.UserId, req.RouteId, req.Comment, int(req.Rating))
	if err != nil {
		return nil, err
	}
	return &pb.Review{
		Id:        r.ID,
		UserId:    r.UserID,
		RouteId:   r.RouteID,
		Rating:    int32(r.Rating),
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
	}, nil
}

func main() {
	_ = godotenv.Load()

	// 1️⃣ DB + migración
	conn, err := reviewsdb.Connect()
	if err != nil {
		log.Fatalf("[reviews] ❌ DB error: %v", err)
	}

	if err := conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("[reviews] ❌ UUID extension error: %v", err)
	}

	if err := conn.AutoMigrate(&model.Review{}); err != nil {
		log.Fatalf("[reviews] ❌ migration error: %v", err)
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
	pb.RegisterReviewsServer(grpcServer, &reviewServer{ctrl: ctrl})

	// Health gRPC
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 3️⃣ Health HTTP (para probes)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK - reviews")
		})
		log.Printf("[reviews] health HTTP on :%d", healthHTTPPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", healthHTTPPort), nil); err != nil {
			log.Printf("[reviews] health HTTP error: %v", err)
		}
	}()

	log.Printf("[reviews] readiness HTTP on :%d", healthHTTPPort)

	// 4️⃣ Servidor principal
	go func() {
		log.Printf("[reviews] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[reviews] server error: %v", err)
		}
	}()

	// 5️⃣ Apagado elegante
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

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
