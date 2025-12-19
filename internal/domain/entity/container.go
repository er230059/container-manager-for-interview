package entity

import "github.com/moby/moby/api/types/container"

type Container struct {
	ID     string
	Image  string
	Cmd    []string
	Env    []string
	Status container.ContainerState
}

func NewContainer(id string, userId int64, image string, cmd []string, env []string, status container.ContainerState) *Container {
	return &Container{
		ID:     id,
		Image:  image,
		Cmd:    cmd,
		Env:    env,
		Status: status,
	}
}
