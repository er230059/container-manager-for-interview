package handler

import (
	"bytes"
	"container-manager/internal/application"
	"container-manager/internal/application/mocks"
	"container-manager/internal/domain/entity"
	"container-manager/internal/errors"
	"container-manager/internal/server/middleware"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	idNode, _ := snowflake.NewNode(1)
	userService := application.NewUserService(mockUserRepo, idNode, "secret")
	userHandler := NewUserHandler(userService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.POST("/users", userHandler.CreateUser)

	t.Run("success", func(t *testing.T) {
		reqBody := CreateUserRequest{
			Username: "testuser",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp CreateUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Username, resp.Username)
		assert.NotEmpty(t, resp.ID)
	})

	t.Run("invalid request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	idNode, _ := snowflake.NewNode(1)
	userService := application.NewUserService(mockUserRepo, idNode, "secret")
	userHandler := NewUserHandler(userService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.POST("/users/login", userHandler.Login)

	t.Run("success", func(t *testing.T) {
		username := "testuser"
		password := "password123"
		
		// Create a real user with hashed password for the mock to return
		user, _ := entity.NewUser(12345, username, password)

		reqBody := LoginRequest{
			Username: username,
			Password: password,
		}
		body, _ := json.Marshal(reqBody)

		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil)

		req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, username, resp.Username)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("user not found", func(t *testing.T) {
		reqBody := LoginRequest{
			Username: "unknown",
			Password: "password",
		}
		body, _ := json.Marshal(reqBody)

		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), "unknown").Return(nil, errors.UserNotFound)

		req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Handler returns Unauthorized (401) when user not found or password mismatch
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}