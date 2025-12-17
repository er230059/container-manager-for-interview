package repository

import (
	"container-manager/internal/domain/entity"
	"context"
)

type ContainerRepository interface {
	Create(ctx context.Context, container *entity.Container) error
	Delete(ctx context.Context, id string) error
}
