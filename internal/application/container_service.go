package application

import (
	"context"

	containerruntime "container-manager/internal/domain"
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
)

type ContainerService struct {
	runtime    containerruntime.ContainerRuntime
	repository repository.ContainerRepository
}

func NewContainerService(runtime containerruntime.ContainerRuntime, repository repository.ContainerRepository) *ContainerService {
	return &ContainerService{runtime: runtime, repository: repository}
}

func (s *ContainerService) CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (string, error) {
	id, err := s.runtime.Create(ctx, options)
	if err != nil {
		return "", err
	}

	container := entity.NewContainer(id, userID, options.Image)

	err = s.repository.Create(ctx, container)
	if err != nil {
		// Here we might want to handle the case where the container was created in the runtime
		// but we failed to save it to the database. For now, we just return the error.
		return "", err
	}

	return id, nil
}

func (s *ContainerService) StartContainer(ctx context.Context, id string) error {
	return s.runtime.Start(ctx, id)
}

func (s *ContainerService) StopContainer(ctx context.Context, id string) error {
	return s.runtime.Stop(ctx, id)
}
