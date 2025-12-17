package repository

import (
	"context"
)

type ContainerUserRepository interface {
	Create(ctx context.Context, containerID string, userID int64) error
	Delete(ctx context.Context, containerID string) error
}
