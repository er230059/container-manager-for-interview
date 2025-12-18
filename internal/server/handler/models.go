package handler

import (
	"encoding/json"
	"time"
)

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

type JobIDResponse struct {
	JobID string `json:"job_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateContainerRequest struct {
	Cmd   []string `json:"cmd"`
	Env   []string `json:"env"`
	Image string   `json:"image" binding:"required"`
}

type JobStatusResponse struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
