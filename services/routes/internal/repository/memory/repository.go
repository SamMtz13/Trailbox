package memory

import (
	"fmt"
	"sync"

	"trailbox/services/routes/internal/model"
	"trailbox/services/routes/internal/repository"
)

type memoryRepository struct {
	mu     sync.RWMutex
	routes map[int]*model.Route
}

// constructor
func NewRepository() repository.Repository {
	return &memoryRepository{routes: make(map[int]*model.Route)}
}

func (r *memoryRepository) Create(route *model.Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[route.ID] = route
	return nil
}

func (r *memoryRepository) GetByID(id int) (*model.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if route, ok := r.routes[id]; ok {
		return route, nil
	}
	return nil, fmt.Errorf("route not found")
}

func (r *memoryRepository) List() ([]*model.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := []*model.Route{}
	for _, route := range r.routes {
		result = append(result, route)
	}
	return result, nil
}
