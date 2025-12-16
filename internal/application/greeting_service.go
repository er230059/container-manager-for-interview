package application

import "leadtek/internal/domain"

// GreetingService is the application service that orchestrates domain logic.
type GreetingService struct {
	repo domain.GreetingRepository
}

// NewGreetingService creates a new GreetingService.
func NewGreetingService(repo domain.GreetingRepository) *GreetingService {
	return &GreetingService{repo: repo}
}

// GetGreeting retrieves the greeting message.
func (s *GreetingService) GetGreeting() (string, error) {
	g, err := s.repo.GetGreeting()
	if err != nil {
		return "", err
	}
	return g.Message, nil
}
