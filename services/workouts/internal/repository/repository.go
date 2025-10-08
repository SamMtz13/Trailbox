package repository

import (
	"trailbox/services/workouts/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	Create(w *model.Workout) error
	GetByID(id uuid.UUID) (*model.Workout, error)
	List() ([]*model.Workout, error)
}
