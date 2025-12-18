package handler

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ContainerIDResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateContainerRequest struct {
	Cmd   []string `json:"cmd"`
	Env   []string `json:"env"`
	Image string   `json:"image" binding:"required"`
}
