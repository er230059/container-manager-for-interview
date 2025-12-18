package repository

import (
	"context"
	"errors"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
	"container-manager/internal/infrastructure/database"
)

var _ repository.ContainerRepository = (*containerRepository)(nil)

type containerRepository struct {
	runtime                 containerruntime.ContainerRuntime
	containerUserRepository database.ContainerUser
}

func NewContainerRepository(
	runtime containerruntime.ContainerRuntime,
	containerUserRepository database.ContainerUser,
) repository.ContainerRepository {
	return &containerRepository{
		runtime:                 runtime,
		containerUserRepository: containerUserRepository,
	}
}

func (r *containerRepository) CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (string, error) {
	id, err := r.runtime.Create(ctx, options)
	if err != nil {
		return "", err
	}

	container := entity.NewContainer(id, userID, options.Image)

	err = r.containerUserRepository.Create(ctx, container.ID, container.UserID)
	if err != nil {
		// TODO: Here we might want to handle the case where the container was created in the runtime
		// but we failed to save it to the database. For now, we just return the error.
		return "", err
	}

	return id, nil
}

func (r *containerRepository) StartContainer(ctx context.Context, userID int64, id string) error {
	ownerID, err := r.containerUserRepository.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("permission denied")
	}

	return r.runtime.Start(ctx, id)
}

func (r *containerRepository) StopContainer(ctx context.Context, userID int64, id string) error {
	ownerID, err := r.containerUserRepository.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("permission denied")
	}

	return r.runtime.Stop(ctx, id)
}

func (r *containerRepository) RemoveContainer(ctx context.Context, userID int64, id string) error {
	ownerID, err := r.containerUserRepository.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return err
	}

	if ownerID != userID {
		return errors.New("permission denied")
	}

	if err := r.runtime.Remove(ctx, id); err != nil {
		return err
	}
	return r.containerUserRepository.Delete(ctx, id)
}
