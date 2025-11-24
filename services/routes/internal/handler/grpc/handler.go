package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonpb "trailbox/gen/common"
	pb "trailbox/gen/routes"
	routesctrl "trailbox/services/routes/internal/controller/routes"
)

type Handler struct {
	pb.UnimplementedRoutesServer
	ctrl *routesctrl.Controller
}

func New(ctrl *routesctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetRoute(ctx context.Context, req *commonpb.RouteId) (*pb.Route, error) {
	route, err := h.ctrl.GetRoute(req.Id)
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

func (h *Handler) ListRoutes(ctx context.Context, req *pb.ListRoutesRequest) (*pb.ListRoutesResponse, error) {
	routes, err := h.ctrl.ListRoutes()
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
