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

	pb "trailbox/gen/reviews"
	reviewsctrl "trailbox/services/reviews/internal/controller"
	consulreg "trailbox/services/reviews/internal/discovery/consul"
	reviewsdb "trailbox/services/reviews/internal/repository/db"
)

const defaultPort = "50051"

type reviewServer struct {
	pb.UnimplementedReviewsServer
	ctrl *reviewsctrl.Controller
}

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
	port := getenvOr("PORT", defaultPort)

	conn, err := reviewsdb.Connect()
	if err != nil {
		log.Fatalf("[reviews] DB error: %v", err)
	}

	repo := reviewsdb.New(conn)
	ctrl := reviewsctrl.NewController(repo)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[reviews] failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterReviewsServer(s, &reviewServer{ctrl: ctrl})

	hs := health.NewServer()
	healthpb.RegisterHealthServer(s, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reg, err := consulreg.NewRegistrar()
	if err != nil {
		log.Fatalf("[reviews] consul init error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "reviews")
	healthPath := getenvOr("SERVICE_HEALTH_PATH", "/grpc.health.v1.Health/Check")
	id, err := reg.Register(getenvOr("SERVICE_NAME", "reviews"), addr, mustAtoi(port), healthPath)
	if err != nil {
		log.Fatalf("[reviews] consul register error: %v", err)
	}
	log.Printf("[reviews] registered in consul id=%s", id)

	go func() {
		log.Printf("[reviews] listening on :%s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[reviews] server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("[reviews] shutting down...")
	s.GracefulStop()
	reg.Deregister()
	log.Println("[reviews] shutdown complete")
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
