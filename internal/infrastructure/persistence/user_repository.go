package persistence

import (
	"fmt"
	"leadtek/internal/domain"
	"sync"
)

// InmemUserRepository is an in-memory implementation of the UserRepository.
type InmemUserRepository struct {
	mu    sync.RWMutex
	users map[int64]*domain.User
}

// NewInmemUserRepository creates a new in-memory user repository.
func NewInmemUserRepository() domain.UserRepository {
	return &InmemUserRepository{
		users: make(map[int64]*domain.User),
	}
}

// Create saves a new user to the in-memory store.
func (r *InmemUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// In a real repo, you'd check for username uniqueness.
	if _, ok := r.users[user.ID]; ok {
		return fmt.Errorf("user with id %d already exists", user.ID)
	}

	r.users[user.ID] = user
	return nil
}
