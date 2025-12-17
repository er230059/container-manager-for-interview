package application

import (
	"context"
	"errors"

	containerruntime "container-manager/internal/domain"
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
)

type ContainerService struct {
	runtime                 containerruntime.ContainerRuntime
	containerUserRepository repository.ContainerUserRepository
}

func NewContainerService(runtime containerruntime.ContainerRuntime, containerUserRepository repository.ContainerUserRepository) *ContainerService {
	return &ContainerService{runtime: runtime, containerUserRepository: containerUserRepository}
}

func (s *ContainerService) CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (string, error) {
	id, err := s.runtime.Create(ctx, options)
	if err != nil {
		return "", err
	}

	container := entity.NewContainer(id, userID, options.Image)

	err = s.containerUserRepository.Create(ctx, container.ID, container.UserID)
	if err != nil {
		// TODO: Here we might want to handle the case where the container was created in the runtime
		// but we failed to save it to the database. For now, we just return the error.
		return "", err
	}

	return id, nil
}

func (s *ContainerService) StartContainer(ctx context.Context, userID int64, id string) error {
	ownerID, err := s.containerUserRepository.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("permission denied")
	}

	return s.runtime.Start(ctx, id)
}

func (s *ContainerService) StopContainer(ctx context.Context, userID int64, id string) error {
	ownerID, err := s.containerUserRepository.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("permission denied")
	}

	return s.runtime.Stop(ctx, id)
}

func (s *ContainerService) RemoveContainer(ctx context.Context, userID int64, id string) error {
	ownerID, err := s.containerUserRepository.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("permission denied")
	}

	if err := s.runtime.Remove(ctx, id); err != nil {
		return err
	}
	return s.containerUserRepository.Delete(ctx, id)
}
