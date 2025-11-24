package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonpb "trailbox/gen/common"
	pb "trailbox/gen/workouts"
	wctrl "trailbox/services/workouts/internal/controller/workouts"
)

type Handler struct {
	pb.UnimplementedWorkoutsServer
	ctrl *wctrl.Controller
}

func New(ctrl *wctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetWorkout(ctx context.Context, req *commonpb.UserId) (*pb.Workout, error) {
	w, err := h.ctrl.GetWorkout(req.Id)
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

func (h *Handler) ListWorkouts(ctx context.Context, req *pb.ListWorkoutsRequest) (*pb.ListWorkoutsResponse, error) {
	workouts, err := h.ctrl.ListWorkouts()
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
