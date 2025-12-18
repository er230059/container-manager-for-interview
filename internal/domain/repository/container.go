package repository

import (
	"context"

	"container-manager/internal/domain/entity"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
)

type ContainerRepository interface {
	CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (*entity.Container, error)
	StartContainer(ctx context.Context, userID int64, id string) error
	StopContainer(ctx context.Context, userID int64, id string) error
	RemoveContainer(ctx context.Context, userID int64, id string) error
	GetContainer(ctx context.Context, id string) (*entity.Container, error)
}
