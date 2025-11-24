package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonpb "trailbox/gen/common"
	pb "trailbox/gen/users"
	userctrl "trailbox/services/users/internal/controller/users"
)

type Handler struct {
	pb.UnimplementedUsersServer
	ctrl *userctrl.Controller
}

func New(ctrl *userctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetUser(ctx context.Context, req *commonpb.UserId) (*pb.User, error) {
	user, err := h.ctrl.GetUser(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &pb.User{
		Id:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h *Handler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := h.ctrl.ListUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}
	resp := &pb.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &pb.User{
			Id:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return resp, nil
}
