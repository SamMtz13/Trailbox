package repository

import (
	"context"
	"trailbox/services/users/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetUser(id string) (*model.User, error)
	ListUsers() ([]model.User, error)
}
