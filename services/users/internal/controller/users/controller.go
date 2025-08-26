package users

import (
	"trailbox/users/internal/model"
	"trailbox/users/internal/repository"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) AddUser(id, name string) error {
	u := &model.User{ID: id, Name: name}
	return c.repo.Create(u)
}

func (c *Controller) GetUser(id string) (*model.User, error) {
	return c.repo.GetByID(id)
}

func (c *Controller) ListUsers() ([]*model.User, error) {
	return c.repo.List()
}
