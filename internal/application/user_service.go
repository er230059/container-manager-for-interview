package application

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo  repository.UserRepository
	idNode    *snowflake.Node
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository, idNode *snowflake.Node, jwtSecret string) *UserService {
	return &UserService{userRepo: userRepo, idNode: idNode, jwtSecret: jwtSecret}
}

func (s *UserService) CreateUser(ctx context.Context, username, plainPassword string) (*entity.User, error) {
	id := s.idNode.Generate().Int64()

	user, err := entity.NewUser(id, username, plainPassword)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*entity.User, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
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
		"sub": strconv.FormatInt(user.ID, 10),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err := claims.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
