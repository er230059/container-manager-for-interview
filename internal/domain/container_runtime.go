package containerruntime

import (
	"context"
)

type ContainerCreateOptions struct {
	Image string
	Env   []string
}

type ContainerRuntime interface {
	Create(ctx context.Context, options ContainerCreateOptions) (string, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
}
