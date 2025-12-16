package handler

import (
	"container-manager/internal/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related API requests.
type UserHandler struct {
	service *application.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// createUserRequest is the request body for creating a user.
type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateUser handles the POST /users request.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(req.Username, req.Password)
	if err != nil {
		// In a real app, you'd check for specific error types,
		// e.g., username already exists, and return a 409 Conflict.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Return the created user, but without the password.
	c.JSON(http.StatusCreated, gin.H{
		"id":       strconv.FormatInt(user.ID, 10),
		"username": user.Username,
	})
}
