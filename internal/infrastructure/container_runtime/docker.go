package containerruntime

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/infrastructure"
	"context"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

var _ infrastructure.ContainerRuntime = (*DockerContainerRuntime)(nil)

type DockerContainerRuntime struct {
	client *client.Client
}

func NewDockerContainerRuntime() (*DockerContainerRuntime, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerContainerRuntime{client: cli}, nil
}

func (d *DockerContainerRuntime) Create(ctx context.Context, options infrastructure.ContainerCreateOptions) (string, error) {
	out, err := d.client.ImagePull(ctx, options.Image, client.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(io.Discard, out)

	resp, err := d.client.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				Cmd:   options.Cmd,
				Env:   options.Env,
				Image: options.Image,
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (d *DockerContainerRuntime) Start(ctx context.Context, id string) error {
	_, err := d.client.ContainerStart(ctx, id, client.ContainerStartOptions{})
	return err
}

func (d *DockerContainerRuntime) Stop(ctx context.Context, id string) error {
	_, err := d.client.ContainerStop(ctx, id, client.ContainerStopOptions{})
	return err
}

func (d *DockerContainerRuntime) Remove(ctx context.Context, id string) error {
	// Compiler said ContainerRemove returns (ContainerRemoveResult, error)
	_, err := d.client.ContainerRemove(ctx, id, client.ContainerRemoveOptions{Force: true})
	return err
}

func (d *DockerContainerRuntime) Inspect(ctx context.Context, id string) (*entity.Container, error) {
	resp, err := d.client.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return nil, err
	}
	// The structure in original code was resp.Container.Config...
	// Let's verify if ContainerInspectResult has Container field.
	// Error msg said: resp.Config undefined (type client.ContainerInspectResult has no field or method Config)
	// So it likely has Container or something else.
	// Based on original code:
	// return &entity.Container{
	// 	ID:     id,
	// 	Image:  resp.Container.Config.Image,
	// 	...
	// }
	// This implies resp has Container field.

	return &entity.Container{
		ID:     id,
		Image:  resp.Container.Config.Image,
		Cmd:    resp.Container.Config.Cmd,
		Env:    resp.Container.Config.Env,
		Status: resp.Container.State.Status,
	}, nil
}
