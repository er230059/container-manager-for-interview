package application

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
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

func (s *UserService) Login(username, password string) (*entity.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errors.New("user not found")
	}

	err = user.ValidatePassword(password)
	if err != nil {
		return nil, "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	// It's better to use a secret from config
	token, err := claims.SignedString([]byte("your-secret-key"))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
