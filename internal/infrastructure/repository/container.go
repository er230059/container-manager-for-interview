package repository

import (
	"context"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	containerruntime "container-manager/internal/infrastructure/container_runtime"
	"container-manager/internal/infrastructure/database"
)

var _ repository.ContainerRepository = (*containerRepository)(nil)

type containerRepository struct {
	runtime               containerruntime.ContainerRuntime
	containerUserDatabase database.ContainerUserDatabase
}

func NewContainerRepository(
	runtime containerruntime.ContainerRuntime,
	containerUserRepository database.ContainerUserDatabase,
) repository.ContainerRepository {
	return &containerRepository{
		runtime:               runtime,
		containerUserDatabase: containerUserRepository,
	}
}

func (r *containerRepository) CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (*entity.Container, error) {
	id, err := r.runtime.Create(ctx, options)
	if err != nil {
		return nil, err
	}

	container := entity.NewContainer(id, userID, options.Image, options.Cmd, options.Env)

	err = r.containerUserDatabase.Create(ctx, container.ID, container.UserID)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (r *containerRepository) StartContainer(ctx context.Context, userID int64, id string) error {
	return r.runtime.Start(ctx, id)
}

func (r *containerRepository) StopContainer(ctx context.Context, userID int64, id string) error {
	return r.runtime.Stop(ctx, id)
}

func (r *containerRepository) RemoveContainer(ctx context.Context, userID int64, id string) error {
	if err := r.runtime.Remove(ctx, id); err != nil {
		return err
	}
	return r.containerUserDatabase.Delete(ctx, id)
}

func (r *containerRepository) GetContainer(ctx context.Context, id string) (*entity.Container, error) {
	ownerID, err := r.containerUserDatabase.GetUserIDByContainerID(ctx, id)
	if err != nil {
		return nil, err
	}
	info, err := r.runtime.Inspect(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity.NewContainer(id, ownerID, info.Image, info.Cmd, info.Env), nil
}
