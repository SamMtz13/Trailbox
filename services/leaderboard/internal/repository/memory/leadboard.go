package memory

import (
	"context"
	"sync"

	"trailbox/services/leaderboard/internal/controller"
)

type repo struct {
	mu   sync.RWMutex
	data map[string]controller.Entry
}

func NewRepository() controller.Repository {
	return &repo{data: make(map[string]controller.Entry)}
}

func (r *repo) Upsert(ctx context.Context, e controller.Entry) (controller.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[e.UserID] = e
	return e, nil
}

func (r *repo) Top(ctx context.Context, limit int) ([]controller.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]controller.Entry, 0, len(r.data))
	for _, v := range r.data {
		out = append(out, v)
	}
	return out, nil
}
