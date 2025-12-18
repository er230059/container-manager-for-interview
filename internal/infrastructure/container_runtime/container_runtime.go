package containerruntime

import (
	"context"
)

type ContainerCreateOptions struct {
	Cmd   []string
	Env   []string
	Image string
}

type ContainerRuntime interface {
	Create(ctx context.Context, options ContainerCreateOptions) (string, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
}
