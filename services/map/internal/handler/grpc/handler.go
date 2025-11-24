package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "trailbox/gen/maps"
	mapctrl "trailbox/services/map/internal/controller"
)

type Handler struct {
	pb.UnimplementedMapServer
	ctrl *mapctrl.Controller
}

func New(ctrl *mapctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetRoute(ctx context.Context, req *pb.GetRouteRequest) (*pb.GetRouteResponse, error) {
	m, err := h.ctrl.GetRouteMap(req.RouteId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "route map not found")
	}
	return &pb.GetRouteResponse{
		RouteId:   m.RouteID.String(),
		GeoJson:   m.GeoJSON,
		CreatedAt: m.CreatedAt.String(),
	}, nil
}

func (h *Handler) SetRoute(ctx context.Context, req *pb.SetRouteRequest) (*pb.SetRouteResponse, error) {
	if err := h.ctrl.SetRouteMap(req.RouteId, req.GeoJson); err != nil {
		return nil, status.Error(codes.Internal, "failed to save map")
	}
	return &pb.SetRouteResponse{Ok: true}, nil
}
