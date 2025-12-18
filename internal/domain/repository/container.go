package repository

import (
	"context"

	containerruntime "container-manager/internal/infrastructure/container_runtime"
)

type ContainerRepository interface {
	CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (string, error)
	StartContainer(ctx context.Context, userID int64, id string) error
	StopContainer(ctx context.Context, userID int64, id string) error
	RemoveContainer(ctx context.Context, userID int64, id string) error
}