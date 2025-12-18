package containerruntime

import (
	"container-manager/internal/domain/infrastructure"
	"context"
)

type ContainerRuntime interface {
	Create(ctx context.Context, options *infrastructure.ContainerCreateOptions) (string, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Inspect(ctx context.Context, id string) (*infrastructure.ContainerInfo, error)
}
