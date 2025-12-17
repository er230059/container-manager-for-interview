package application

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"

	"github.com/bwmarrin/snowflake"
)

type UserService struct {
	userRepo repository.UserRepository
	idNode   *snowflake.Node
}

func NewUserService(userRepo repository.UserRepository, idNode *snowflake.Node) *UserService {
	return &UserService{userRepo: userRepo, idNode: idNode}
}

func (s *UserService) CreateUser(username, plainPassword string) (*entity.User, error) {
	id := s.idNode.Generate().Int64()

	user, err := entity.NewUser(id, username, plainPassword)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
