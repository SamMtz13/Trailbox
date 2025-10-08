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
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	pb "trailbox/gen/leaderboard"

	lbctrl "trailbox/services/leaderboard/internal/controller"
	lbdb "trailbox/services/leaderboard/internal/db"
	lbconsul "trailbox/services/leaderboard/internal/discovery/consul"
	lbrepo "trailbox/services/leaderboard/internal/repository/db"
)

type leaderboardServer struct {
	pb.UnimplementedLeaderboardServer
	ctrl *lbctrl.Controller
}

func (s *leaderboardServer) GetTop(ctx context.Context, req *pb.GetTopRequest) (*pb.GetTopResponse, error) {
	rows, err := s.ctrl.GetTop(int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get top")
	}
	resp := &pb.GetTopResponse{}
	for _, r := range rows {
		resp.Entries = append(resp.Entries, &pb.LeaderboardEntry{
			Id:       r.ID.String(),
			UserId:   r.UserID.String(),
			Score:    int32(r.Score),
			Position: int32(r.Position),
		})
	}
	return resp, nil
}

func (s *leaderboardServer) Upsert(ctx context.Context, req *pb.UpsertRequest) (*pb.UpsertResponse, error) {
	if err := s.ctrl.Upsert(req.UserId, int(req.Score)); err != nil {
		return nil, status.Error(codes.Internal, "failed to upsert score")
	}
	return &pb.UpsertResponse{Ok: true}, nil
}

const defaultPort = "50051"

func main() {
	conn, err := lbdb.Connect()
	if err != nil {
		log.Fatalf("[leaderboard] DB error: %v", err)
	}

	repo := lbrepo.New(conn)
	ctrl := lbctrl.NewController(repo)

	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[leaderboard] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLeaderboardServer(grpcServer, &leaderboardServer{ctrl: ctrl})

	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reg, err := lbconsul.NewRegistrar()
	if err != nil {
		log.Fatalf("[leaderboard] consul error: %v", err)
	}
	addr := getenvOr("SERVICE_ADDRESS", "leaderboard")
	id, err := reg.Register(getenvOr("SERVICE_NAME", "leaderboard"), addr, mustAtoi(port), "/grpc.health.v1.Health/Check")
	if err != nil {
		log.Fatalf("[leaderboard] consul register error: %v", err)
	}
	log.Printf("[leaderboard] registered in Consul as id=%s", id)

	go func() {
		log.Printf("[leaderboard] listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[leaderboard] serve error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[leaderboard] shutting down...")
	grpcServer.GracefulStop()
	reg.Deregister()
	log.Println("[leaderboard] graceful shutdown complete")
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustAtoi(s string) int { n, _ := strconv.Atoi(s); return n }
