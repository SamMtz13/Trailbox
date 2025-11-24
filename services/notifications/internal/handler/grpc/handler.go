package grpc

import (
	"context"
	"time"

	pb "trailbox/gen/notifications"
	notifctrl "trailbox/services/notifications/internal/controller"
)

type Handler struct {
	pb.UnimplementedNotificationsServer
	ctrl *notifctrl.Controller
}

func New(ctrl *notifctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetNotifications(ctx context.Context, req *pb.UserIdRequest) (*pb.NotificationsResponse, error) {
	notifs, err := h.ctrl.ListByUser(req.UserId)
	if err != nil {
		return nil, err
	}

	resp := &pb.NotificationsResponse{}
	for _, n := range notifs {
		resp.Notifications = append(resp.Notifications, &pb.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			Message:   n.Message,
			Read:      n.Read,
			CreatedAt: n.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

func (h *Handler) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.Notification, error) {
	n, err := h.ctrl.Create(req.UserId, req.Message)
	if err != nil {
		return nil, err
	}
	return &pb.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Message:   n.Message,
		Read:      n.Read,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}, nil
}
