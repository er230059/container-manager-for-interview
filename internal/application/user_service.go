package application

import (
	"leadtek/internal/domain"

	"github.com/bwmarrin/snowflake"
)

// UserService handles the application logic for users.
type UserService struct {
	repo   domain.UserRepository
	idNode *snowflake.Node
}

// NewUserService creates a new UserService.
func NewUserService(repo domain.UserRepository, idNode *snowflake.Node) *UserService {
	return &UserService{repo: repo, idNode: idNode}
}

// CreateUser creates a new user, generates an ID, and persists it.
func (s *UserService) CreateUser(username, password string) (*domain.User, error) {
	user := &domain.User{
		ID:       s.idNode.Generate().Int64(),
		Username: username,
		Password: password,
	}

	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
