package repository

import (
	"context"

	"trailbox/services/routes/internal/model"
)

type Repository interface {
	CreateRoute(ctx context.Context, route *model.Route) error
	GetRoute(id string) (*model.Route, error)
	ListRoutes() ([]model.Route, error)
}
