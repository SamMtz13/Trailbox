package repository

import "trailbox/services/users/internal/model"

type Repository interface {
	Create(u *model.User) error
	GetByID(id string) (*model.User, error)
	List() ([]*model.User, error)
}
