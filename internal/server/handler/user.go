package handler

// @BasePath /

import (
	"container-manager/internal/application"
	"container-manager/internal/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User creation request"
// @Success 200 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err))
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, CreateUserResponse{
		ID:       strconv.FormatInt(user.ID, 10),
		Username: user.Username,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns an authentication token
// @Tags Users
// @Accept json
// @Produce json
// @Param user body LoginRequest true "User login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err))
		return
	}

	user, token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" || err.Error() == "user not found" {
			_ = c.Error(errors.Unauthorized)
			return
		}
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		ID:       strconv.FormatInt(user.ID, 10),
		Username: user.Username,
		Token:    token,
	})
}
