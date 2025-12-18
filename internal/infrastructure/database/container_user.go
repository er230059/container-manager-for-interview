package database

import "context"

type ContainerUser interface {
	Create(ctx context.Context, containerID string, userID int64) error
	Delete(ctx context.Context, containerID string) error
	GetUserIDByContainerID(ctx context.Context, containerID string) (int64, error)
}
