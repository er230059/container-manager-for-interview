package containerruntime

import (
	"context"
	"io"
	"os"

	containerruntime "container-manager/internal/domain"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

var _ containerruntime.ContainerRuntime = (*DockerContainerRuntime)(nil)

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

func (d *DockerContainerRuntime) Create(ctx context.Context, options containerruntime.ContainerCreateOptions) (string, error) {
	out, err := d.client.ImagePull(ctx, options.Image, client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := d.client.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				Image: options.Image,
				Env:   options.Env,
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
	_, err := d.client.ContainerRemove(ctx, id, client.ContainerRemoveOptions{})
	return err
}
