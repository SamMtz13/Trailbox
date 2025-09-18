package controller

import (
	"context"
	"time"
)

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Channel   string    `json:"channel"` // email|push|sms
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Store(ctx context.Context, n Notification) (Notification, error)
	List(ctx context.Context) ([]Notification, error)
}

type Controller struct{ repo Repository }

func NewController(r Repository) *Controller { return &Controller{repo: r} }

func (c *Controller) Send(ctx context.Context, n Notification) (Notification, error) {
	n.CreatedAt = time.Now().UTC()
	n.Status = "sent" // simulamos que siempre se envía bien
	return c.repo.Store(ctx, n)
}

func (c *Controller) List(ctx context.Context) ([]Notification, error) {
	return c.repo.List(ctx)
}
