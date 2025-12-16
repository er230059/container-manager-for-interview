package persistence

import "leadtek/internal/domain"

// inmemGreetingRepo is an in-memory implementation of the GreetingRepository.
type inmemGreetingRepo struct{}

// NewInmemGreetingRepository creates a new in-memory repository.
func NewInmemGreetingRepository() domain.GreetingRepository {
	return &inmemGreetingRepo{}
}

// GetGreeting fetches a greeting from the in-memory store.
func (r *inmemGreetingRepo) GetGreeting() (*domain.Greeting, error) {
	// In a real application, this would fetch from a database.
	return &domain.Greeting{Message: "This is home page"}, nil
}
