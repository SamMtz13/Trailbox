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
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	pb "trailbox/gen/maps"

	mapctrl "trailbox/services/map/internal/controller"
	mapdb "trailbox/services/map/internal/db"
	maprepo "trailbox/services/map/internal/repository/db"
)

const (
	defaultPort    = "50051"
	healthHTTPPort = 8081
)

type mapServer struct {
	pb.UnimplementedMapServer
	ctrl *mapctrl.Controller
}

func (s *mapServer) GetRoute(ctx context.Context, req *pb.GetRouteRequest) (*pb.GetRouteResponse, error) {
	m, err := s.ctrl.GetRouteMap(req.RouteId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "route map not found")
	}
	return &pb.GetRouteResponse{
		RouteId:   m.RouteID.String(),
		GeoJson:   m.GeoJSON,
		CreatedAt: m.CreatedAt.String(),
	}, nil
}

func (s *mapServer) SetRoute(ctx context.Context, req *pb.SetRouteRequest) (*pb.SetRouteResponse, error) {
	if err := s.ctrl.SetRouteMap(req.RouteId, req.GeoJson); err != nil {
		return nil, status.Error(codes.Internal, "failed to save map")
	}
	return &pb.SetRouteResponse{Ok: true}, nil
}

func main() {
	conn, err := mapdb.Connect()
	if err != nil {
		log.Fatalf("[map] DB error: %v", err)
	}

	repo := maprepo.New(conn)
	ctrl := mapctrl.NewController(repo)

	port := getenvOr("PORT", defaultPort)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[map] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMapServer(grpcServer, &mapServer{ctrl: ctrl})

	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// HTTP health para probes
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK - map")
		})
		log.Printf("[map] HTTP health running on :%d", healthHTTPPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", healthHTTPPort), nil); err != nil {
			log.Printf("[map] health server error: %v", err)
		}
	}()

	log.Printf("[map] readiness HTTP on :%d", healthHTTPPort)

	go func() {
		log.Printf("[map] 🚀 gRPC listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[map] serve error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[map] shutting down...")
	grpcServer.GracefulStop()
	log.Println("[map] graceful shutdown complete")
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustAtoi(s string) int { n, _ := strconv.Atoi(s); return n }
