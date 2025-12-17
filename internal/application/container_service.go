package application

import (
	"context"

	containerruntime "container-manager/internal/domain"
)

type ContainerService struct {
	runtime containerruntime.ContainerRuntime
}

func NewContainerService(runtime containerruntime.ContainerRuntime) *ContainerService {
	return &ContainerService{runtime: runtime}
}

func (s *ContainerService) CreateContainer(ctx context.Context, options containerruntime.ContainerCreateOptions) (string, error) {
	return s.runtime.Create(ctx, options)
}
