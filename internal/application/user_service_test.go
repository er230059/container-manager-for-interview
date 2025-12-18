package application

import (
	"container-manager/internal/application/mocks"
	"container-manager/internal/domain/entity"
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestUserService_CreateUser(t *testing.T) {
	username := "testuser"
	password := "testpassword"

	t.Run("successful user creation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		idNode, _ := snowflake.NewNode(1)
		userService := NewUserService(mockUserRepo, idNode, "test_secret")

		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		user, err := userService.CreateUser(context.Background(), username, password)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.Username != username {
			t.Errorf("expected username %s, got %s", username, user.Username)
		}
	})

	t.Run("user creation failed by repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		idNode, _ := snowflake.NewNode(1)
		userService := NewUserService(mockUserRepo, idNode, "test_secret")

		expectedErr := errors.New("database error")
		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedErr).Times(1)

		user, err := userService.CreateUser(context.Background(), username, password)
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if user != nil {
			t.Errorf("expected nil user, got %v", user)
		}
	})
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	idNode, _ := snowflake.NewNode(1)
	userService := NewUserService(mockUserRepo, idNode, "test_secret")

	username := "testuser"
	plainPassword := "testpassword"

	t.Run("successful login", func(t *testing.T) {
		user, _ := entity.NewUser(idNode.Generate().Int64(), username, plainPassword)
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil).Times(1)

		loggedInUser, token, err := userService.Login(context.Background(), username, plainPassword)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if loggedInUser == nil {
			t.Fatal("expected loggedInUser, got nil")
		}
		if loggedInUser.Username != username {
			t.Errorf("expected username %s, got %s", username, loggedInUser.Username)
		}
		if token == "" {
			t.Error("expected token, got empty string")
		}

		// Verify JWT token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("test_secret"), nil
		})
		if err != nil {
			t.Fatalf("failed to parse token: %v", err)
		}
		if !parsedToken.Valid {
			t.Error("token is invalid")
		}
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("failed to parse claims")
		}
		if claims["sub"] != strconv.FormatInt(loggedInUser.ID, 10) {
			t.Errorf("expected subject %s, got %s", strconv.FormatInt(loggedInUser.ID, 10), claims["sub"])
		}
		exp, ok := claims["exp"].(float64) // JWT standard expiration time is a number
		if !ok {
			t.Fatal("exp claim not found or not a number")
		}
		if int64(exp) <= time.Now().Unix() {
			t.Error("token expired")
		}
	})

	t.Run("login user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, nil).Times(1)

		loggedInUser, token, err := userService.Login(context.Background(), username, plainPassword)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "user not found" {
			t.Errorf("expected 'user not found' error, got %v", err)
		}
		if loggedInUser != nil {
			t.Errorf("expected nil user, got %v", loggedInUser)
		}
		if token != "" {
			t.Errorf("expected empty token, got %s", token)
		}
	})

	t.Run("login wrong password", func(t *testing.T) {
		user, _ := entity.NewUser(idNode.Generate().Int64(), username, plainPassword)
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil).Times(1)

		loggedInUser, token, err := userService.Login(context.Background(), username, "wrongpassword")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			t.Errorf("expected 'crypto/bcrypt: hashedPassword is not the hash of the given password' error, got %v", err)
		}
		if loggedInUser != nil {
			t.Errorf("expected nil user, got %v", loggedInUser)
		}
		if token != "" {
			t.Errorf("expected empty token, got %s", token)
		}
	})

	t.Run("login FindByUsername returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockUserRepo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, expectedErr).Times(1)

		loggedInUser, token, err := userService.Login(context.Background(), username, plainPassword)
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if loggedInUser != nil {
			t.Errorf("expected nil user, got %v", loggedInUser)
		}
		if token != "" {
			t.Errorf("expected empty token, got %s", token)
		}
	})
}
