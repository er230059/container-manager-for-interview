package handler

import (
	"bytes"
	"container-manager/internal/application"
	"container-manager/internal/application/mocks"
	"container-manager/internal/domain/entity"
	"container-manager/internal/server/middleware"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestContainerHandler_ListContainers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	containerService := application.NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)
	containerHandler := NewContainerHandler(containerService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("userID", "123")
		c.Next()
	})
	router.GET("/containers", containerHandler.ListContainers)

	t.Run("success", func(t *testing.T) {
		mockContainerUserRepo.EXPECT().GetContainerIDsByUserID(gomock.Any(), int64(123)).Return([]string{"c1", "c2"}, nil)
		mockRuntime.EXPECT().Inspect(gomock.Any(), "c1").Return(&entity.Container{ID: "c1", Image: "img1", Status: "running"}, nil)
		mockRuntime.EXPECT().Inspect(gomock.Any(), "c2").Return(&entity.Container{ID: "c2", Image: "img2", Status: "stopped"}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/containers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []ContainerResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "c1", resp[0].ID)
		assert.Equal(t, "c2", resp[1].ID)
	})
}

func TestContainerHandler_CreateContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	containerService := application.NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)
	containerHandler := NewContainerHandler(containerService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("userID", "123")
		c.Next()
	})
	router.POST("/containers", containerHandler.CreateContainer)

	t.Run("success", func(t *testing.T) {
		reqBody := CreateContainerRequest{
			Image: "nginx",
			Cmd:   []string{"start"},
			Env:   []string{"ENV=production"},
		}
		body, _ := json.Marshal(reqBody)

		mockJobRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		// The service also starts a goroutine to run the job. 
		// We can't easily assert on that unless we wait or mock the async part, 
		// but `CreateContainer` returns immediately after creating the job.
		// So we just expect `jobRepo.Create`.
		// Note: The goroutine will run and call `jobRepo.Update`, `runtime.Create`, etc.
		// Since those are mocks, they might be called unexpectedly if we don't mock them?
		// But the goroutine runs asynchronously. The test might finish before the goroutine mocks are called.
		// Or it might panic if the controller is finished.
		// To be safe, we can mock the subsequent calls with `.AnyTimes()` or just ignore them 
		// because we are testing the Handler, which only cares about the synchronous part.
		// However, `gomock` Controller detects missing calls or unexpected calls if the goroutine runs fast enough.
		// A common strategy is to allow any calls on mocks for the async part if we don't want to sync.
		// Or mock `jobRepo.Create` to return nil and do nothing else.
		// But the service logic starts the goroutine.
		
		// To avoid race conditions in tests with gomock and goroutines, 
		// we should probably just test the Handler's immediate response.
		// The mocked `jobRepo.Create` is synchronous.
		// The `go s.runCreateContainerJob(...)` starts after.
		// If we are lucky, the test finishes before the goroutine does much.
		// Ideally we should mock the behavior of `runCreateContainerJob` but we can't mock private method.
		
		// Let's allow subsequent calls just in case.
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		mockRuntime.EXPECT().Create(gomock.Any(), gomock.Any()).Return("cid", nil).AnyTimes()
		mockContainerUserRepo.EXPECT().Create(gomock.Any(), "cid", int64(123)).Return(nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodPost, "/containers", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp CreateContainerResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.JobID)
	})
}

func TestContainerHandler_StartContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	containerService := application.NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)
	containerHandler := NewContainerHandler(containerService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("userID", "123")
		c.Next()
	})
	router.PATCH("/containers/:id/start", containerHandler.StartContainer)

	t.Run("success", func(t *testing.T) {
		containerID := "c1"
		mockContainerUserRepo.EXPECT().GetUserIDByContainerID(gomock.Any(), containerID).Return(int64(123), nil)
		mockRuntime.EXPECT().Start(gomock.Any(), containerID).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/containers/c1/start", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("permission denied", func(t *testing.T) {
		containerID := "c2"
		mockContainerUserRepo.EXPECT().GetUserIDByContainerID(gomock.Any(), containerID).Return(int64(456), nil)

		req, _ := http.NewRequest(http.MethodPatch, "/containers/c2/start", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestContainerHandler_StopContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	containerService := application.NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)
	containerHandler := NewContainerHandler(containerService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("userID", "123")
		c.Next()
	})
	router.PATCH("/containers/:id/stop", containerHandler.StopContainer)

	t.Run("success", func(t *testing.T) {
		containerID := "c1"
		mockContainerUserRepo.EXPECT().GetUserIDByContainerID(gomock.Any(), containerID).Return(int64(123), nil)
		mockRuntime.EXPECT().Stop(gomock.Any(), containerID).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/containers/c1/stop", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestContainerHandler_RemoveContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	containerService := application.NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)
	containerHandler := NewContainerHandler(containerService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Set("userID", "123")
		c.Next()
	})
	router.DELETE("/containers/:id", containerHandler.RemoveContainer)

	t.Run("success", func(t *testing.T) {
		containerID := "c1"
		mockContainerUserRepo.EXPECT().GetUserIDByContainerID(gomock.Any(), containerID).Return(int64(123), nil)
		mockRuntime.EXPECT().Remove(gomock.Any(), containerID).Return(nil)
		mockContainerUserRepo.EXPECT().Delete(gomock.Any(), containerID).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/containers/c1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
