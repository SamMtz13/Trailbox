package workouts

import (
	"trailbox/services/workouts/internal/model"
	"trailbox/services/workouts/internal/repository"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}
func (c *Controller) AddWorkout(id, name string) error {
	w := &model.Workout{ID: id, Name: name, Exercises: []string{}, Duration: 0}
	return c.repo.Create(w)
}

func (c *Controller) GetWorkout(id string) (*model.Workout, error) {
	return c.repo.GetByID(id)
}

func (c *Controller) ListWorkouts() ([]*model.Workout, error) {
	return c.repo.List()
}
