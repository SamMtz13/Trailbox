package grpc

import (
	"context"
	"time"

	pb "trailbox/gen/reviews"
	reviewsctrl "trailbox/services/reviews/internal/controller"
)

type Handler struct {
	pb.UnimplementedReviewsServer
	ctrl *reviewsctrl.Controller
}

func New(ctrl *reviewsctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetReviews(ctx context.Context, req *pb.ReviewListRequest) (*pb.ReviewListResponse, error) {
	revs, err := h.ctrl.ListReviews(req.RouteId)
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

func (h *Handler) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.Review, error) {
	r, err := h.ctrl.AddReview(req.UserId, req.RouteId, req.Comment, int(req.Rating))
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
