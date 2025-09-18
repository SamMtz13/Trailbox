package memory

import (
	"context"
	"errors"
	"sync"

	"trailbox/services/map/internal/controller"
)

type repo struct {
	mu   sync.RWMutex
	data map[string]controller.RouteMap
}

func NewRepository() controller.Repository {
	return &repo{data: make(map[string]controller.RouteMap)}
}

func (r *repo) SetRoute(ctx context.Context, rm controller.RouteMap) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[rm.RouteID] = rm
	return nil
}

func (r *repo) GetRoute(ctx context.Context, id string) (controller.RouteMap, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rm, ok := r.data[id]
	if !ok {
		return controller.RouteMap{}, errors.New("not found")
	}
	return rm, nil
}
