package application

import (
	"container-manager/internal/domain"

	"github.com/bwmarrin/snowflake"
)

// UserService handles the application logic for users.
type UserService struct {
	userRepo domain.UserRepository
	idNode   *snowflake.Node
}

// NewUserService creates a new UserService.
func NewUserService(userRepo domain.UserRepository, idNode *snowflake.Node) *UserService {
	return &UserService{userRepo: userRepo, idNode: idNode}
}

// CreateUser creates a new user, generates an ID, and persists it.
func (s *UserService) CreateUser(username, password string) (*domain.User, error) {
	user := &domain.User{
		ID:       s.idNode.Generate().Int64(),
		Username: username,
		Password: password,
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
