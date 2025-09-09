package memory

import (
	"fmt"
	"sync"
	"trailbox/services/workouts/internal/model"
	"trailbox/services/workouts/internal/repository"
)

type memoryRepo struct {
	mu       sync.RWMutex
	workouts map[string]*model.Workout
}

// constructor
func New() repository.Repository {
	return &memoryRepo{workouts: make(map[string]*model.Workout)}
}

func (r *memoryRepo) Create(w *model.Workout) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workouts[w.ID] = w
	return nil
}

func (r *memoryRepo) GetByID(id string) (*model.Workout, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if w, ok := r.workouts[id]; ok {
		return w, nil
	}
	return nil, fmt.Errorf("workout not found")
}

func (r *memoryRepo) List() ([]*model.Workout, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := []*model.Workout{}
	for _, w := range r.workouts {
		result = append(result, w)
	}
	return result, nil
}
