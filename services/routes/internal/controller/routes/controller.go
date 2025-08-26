package routes

import (
	"trailbox/routes/internal/model"
	"trailbox/routes/internal/repository"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) AddRoute(id int, path string) error {
	r := &model.Route{ID: id, Path: path}
	return c.repo.Create(r)
}

func (c *Controller) GetRoute(id int) (*model.Route, error) {
	return c.repo.GetByID(id)
}

func (c *Controller) ListRoutes() ([]*model.Route, error) {
	return c.repo.List()
}

func (c *Controller) CreateRoute(r *model.Route) error {
	return c.repo.Create(r)
}
