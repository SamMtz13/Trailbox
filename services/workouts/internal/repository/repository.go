package repository

import "trailbox/workouts/internal/model"
type Repository interface {
	Create(w *model.Workout) error
	GetByID(id string) (*model.Workout, error)
	List() ([]*model.Workout, error)
}