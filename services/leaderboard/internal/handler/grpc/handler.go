package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "trailbox/gen/leaderboard"
	lbctrl "trailbox/services/leaderboard/internal/controller"
)

type Handler struct {
	pb.UnimplementedLeaderboardServer
	ctrl *lbctrl.Controller
}

func New(ctrl *lbctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetTop(ctx context.Context, req *pb.GetTopRequest) (*pb.GetTopResponse, error) {
	rows, err := h.ctrl.GetTop(int(req.Limit))
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

func (h *Handler) Upsert(ctx context.Context, req *pb.UpsertRequest) (*pb.UpsertResponse, error) {
	if err := h.ctrl.Upsert(req.UserId, int(req.Score)); err != nil {
		return nil, status.Error(codes.Internal, "failed to upsert score")
	}
	return &pb.UpsertResponse{Ok: true}, nil
}
