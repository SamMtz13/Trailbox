package users

import (
	"context"

	"trailbox/services/users/internal/model"
	"trailbox/services/users/internal/repository"

	"github.com/google/uuid"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) AddUser(id, name string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	u := &model.User{ID: uid, Name: name}
	return c.repo.CreateUser(context.TODO(), u) // ✅ reemplazado nil → context.TODO()
}

func (c *Controller) GetUser(id string) (*model.User, error) {
	return c.repo.GetUser(id)
}

func (c *Controller) ListUsers() ([]model.User, error) {
	return c.repo.ListUsers()
}
