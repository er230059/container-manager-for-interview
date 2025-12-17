package entity

type Container struct {
	ID     string
	Image  string
	UserID int64
}

func NewContainer(id string, userId int64, image string) *Container {
	return &Container{
		ID:     id,
		Image:  image,
		UserID: userId,
	}
}
