package controller

import (
	"context"
	"errors"
	"time"
)

type Review struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RouteID   string    `json:"route_id"`
	Rating    int       `json:"rating"` // 1-5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, r Review) (Review, error)
	List(ctx context.Context, routeID, userID string) ([]Review, error)
}

type Controller struct {
	repo Repository
}

func NewController(r Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) Create(ctx context.Context, r Review) (Review, error) {
	if r.UserID == "" || r.RouteID == "" {
		return Review{}, errors.New("user_id and route_id are required")
	}
	if r.Rating < 1 || r.Rating > 5 {
		return Review{}, errors.New("rating must be 1..5")
	}
	r.CreatedAt = time.Now().UTC()
	return c.repo.Create(ctx, r)
}

func (c *Controller) List(ctx context.Context, routeID, userID string) ([]Review, error) {
	return c.repo.List(ctx, routeID, userID)
}
