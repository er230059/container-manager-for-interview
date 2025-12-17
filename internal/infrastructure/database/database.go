package database

import (
	"container-manager/internal/domain/entity"
	"context"
)

type UserDatabase interface {
	CreateUser(ctx context.Context, id int64, username, password string) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
}
