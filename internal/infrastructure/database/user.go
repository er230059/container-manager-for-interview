package database

import (
	"container-manager/internal/domain/entity"
	"context"
)

type User interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
}
