package application

import (
	"context"

	"container-manager/internal/domain/repository"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
)

type ContainerService struct {
	repo repository.ContainerRepository
}

func NewContainerService(repo repository.ContainerRepository) *ContainerService {
	return &ContainerService{repo: repo}
}

func (s *ContainerService) CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (string, error) {
	return s.repo.CreateContainer(ctx, userID, options)
}

func (s *ContainerService) StartContainer(ctx context.Context, userID int64, id string) error {
	return s.repo.StartContainer(ctx, userID, id)
}

func (s *ContainerService) StopContainer(ctx context.Context, userID int64, id string) error {
	return s.repo.StopContainer(ctx, userID, id)
}

func (s *ContainerService) RemoveContainer(ctx context.Context, userID int64, id string) error {
	return s.repo.RemoveContainer(ctx, userID, id)
}
