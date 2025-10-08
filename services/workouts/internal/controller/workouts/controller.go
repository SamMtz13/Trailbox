package workouts

import (
	"trailbox/services/workouts/internal/model"
	"trailbox/services/workouts/internal/repository"

	"github.com/google/uuid"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

// Crear un nuevo workout
func (c *Controller) AddWorkout(name string, exercises []string, duration int, userID, routeID uuid.UUID) error {
	w := &model.Workout{
		Name:      name,
		Exercises: exercises,
		Duration:  duration,
		UserID:    userID,
		RouteID:   routeID,
	}
	return c.repo.Create(w)
}

// Obtener un workout por ID
func (c *Controller) GetWorkout(id string) (*model.Workout, error) {
	workoutID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return c.repo.GetByID(workoutID)
}

// Listar todos los workouts
func (c *Controller) ListWorkouts() ([]*model.Workout, error) {
	return c.repo.List()
}
