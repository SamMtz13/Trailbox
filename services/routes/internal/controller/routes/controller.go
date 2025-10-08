package routes

import (
	"context"

	"trailbox/services/routes/internal/model"
	"trailbox/services/routes/internal/repository"

	"github.com/google/uuid"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) AddRoute(userID, path string, duration, distance int) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	r := &model.Route{
		ID:       uuid.New(),
		Path:     path,
		Duration: duration,
		Distance: distance,
		UserID:   uid,
	}
	return c.repo.CreateRoute(context.TODO(), r)
}

func (c *Controller) GetRoute(id string) (*model.Route, error) {
	return c.repo.GetRoute(id)
}

func (c *Controller) ListRoutes() ([]model.Route, error) {
	return c.repo.ListRoutes()
}
