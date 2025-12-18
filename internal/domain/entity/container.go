package entity

type Container struct {
	ID     string
	Image  string
	Cmd    []string
	Env    []string
	UserID int64
}

func NewContainer(id string, userId int64, image string, cmd []string, env []string) *Container {
	return &Container{
		ID:     id,
		Image:  image,
		Cmd:    cmd,
		Env:    env,
		UserID: userId,
	}
}
