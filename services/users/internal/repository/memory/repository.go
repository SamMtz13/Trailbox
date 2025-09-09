package memory

import (
	"fmt"
	"sync"

	"trailbox/services/users/internal/model"
	"trailbox/services/users/internal/repository"
)

type memoryRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

// constructor
func New() repository.Repository {
	return &memoryRepo{users: make(map[string]*model.User)}
}

func (r *memoryRepo) Create(u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *memoryRepo) GetByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *memoryRepo) List() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := []*model.User{}
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}
