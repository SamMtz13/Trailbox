package repository

import "trailbox/routes/internal/model"

type Repository interface {
	Create(route *model.Route) error
	GetByID(id int) (*model.Route, error)
	List() ([]*model.Route, error)
}
