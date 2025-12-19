package handler

import (
	"encoding/json"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type CreateContainerResponse struct {
	JobID string `json:"job_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateContainerRequest struct {
	Cmd   []string `json:"cmd" example:"tail,-f,/dev/null"`
	Env   []string `json:"env" example:"FOO=BAR"`
	Image string   `json:"image" binding:"required" example:"alpine"`
}

type GetJobResponse struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ContainerResponse struct {
	ID     string   `json:"id"`
	Image  string   `json:"image"`
	Cmd    []string `json:"cmd"`
	Env    []string `json:"env"`
	Status string   `json:"status"`
}
