package memory

import (
	"context"
	"sync"
	"time"

	"trailbox/services/reviews/internal/controller"
)

type repo struct {
	mu      sync.RWMutex
	storage map[string]controller.Review
}

func NewRepository() controller.Repository {
	return &repo{storage: make(map[string]controller.Review)}
}

func genID() string {
	return time.Now().UTC().Format("20060102T150405.000000000")
}

func (r *repo) Create(ctx context.Context, rev controller.Review) (controller.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rev.ID == "" {
		rev.ID = genID()
	}
	r.storage[rev.ID] = rev
	return rev, nil
}

func (r *repo) List(ctx context.Context, routeID, userID string) ([]controller.Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]controller.Review, 0, len(r.storage))
	for _, v := range r.storage {
		if routeID != "" && v.RouteID != routeID {
			continue
		}
		if userID != "" && v.UserID != userID {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}
