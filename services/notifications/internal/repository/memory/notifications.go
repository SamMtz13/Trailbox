package memory

import (
	"context"
	"sync"
	"time"

	"trailbox/services/notifications/internal/controller"
)

type repo struct {
	mu   sync.RWMutex
	data map[string]controller.Notification
}

func NewRepository() controller.Repository {
	return &repo{data: make(map[string]controller.Notification)}
}

func genID() string { return time.Now().UTC().Format("20060102T150405.000000000") }

func (r *repo) Store(ctx context.Context, n controller.Notification) (controller.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n.ID == "" {
		n.ID = genID()
	}
	r.data[n.ID] = n
	return n, nil
}

func (r *repo) List(ctx context.Context) ([]controller.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]controller.Notification, 0, len(r.data))
	for _, v := range r.data {
		out = append(out, v)
	}
	return out, nil
}
