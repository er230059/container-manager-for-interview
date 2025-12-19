package handler

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/errors"
	"container-manager/internal/server/middleware"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockJobService struct {
	GetJobFunc func(ctx context.Context, userID int64, id string) (*entity.Job, error)
}

func (m *MockJobService) GetJob(ctx context.Context, userID int64, id string) (*entity.Job, error) {
	if m.GetJobFunc != nil {
		return m.GetJobFunc(ctx, userID, id)
	}
	return nil, nil
}

func TestJobHandler_GetJob(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &MockJobService{}
		jobHandler := NewJobHandler(mockService)

		router := gin.Default()
		router.Use(middleware.ErrorHandler())
		// Middleware to set userID
		router.Use(func(c *gin.Context) {
			c.Set("userID", "123")
			c.Next()
		})
		router.GET("/jobs/:id", jobHandler.GetJob)

		expectedJob := &entity.Job{
			ID:        "job-1",
			Type:      "test",
			Status:    entity.JobStatusCompleted,
			Result:    []byte(`{"foo":"bar"}`),
			UserID:    123,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.GetJobFunc = func(ctx context.Context, userID int64, id string) (*entity.Job, error) {
			assert.Equal(t, int64(123), userID)
			assert.Equal(t, "job-1", id)
			return expectedJob, nil
		}

		req, _ := http.NewRequest(http.MethodGet, "/jobs/job-1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "job-1")
		assert.Contains(t, w.Body.String(), "completed")
	})

	t.Run("job not found", func(t *testing.T) {
		mockService := &MockJobService{}
		jobHandler := NewJobHandler(mockService)

		router := gin.Default()
		router.Use(middleware.ErrorHandler())
		router.Use(func(c *gin.Context) {
			c.Set("userID", "123")
			c.Next()
		})
		router.GET("/jobs/:id", jobHandler.GetJob)

		mockService.GetJobFunc = func(ctx context.Context, userID int64, id string) (*entity.Job, error) {
			return nil, errors.JobNotFound
		}

		req, _ := http.NewRequest(http.MethodGet, "/jobs/unknown", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockService := &MockJobService{}
		jobHandler := NewJobHandler(mockService)

		router := gin.Default()
		router.Use(middleware.ErrorHandler())
		router.Use(func(c *gin.Context) {
			c.Set("userID", "123")
			c.Next()
		})
		router.GET("/jobs/:id", jobHandler.GetJob)

		mockService.GetJobFunc = func(ctx context.Context, userID int64, id string) (*entity.Job, error) {
			return nil, errors.PermissionDenied
		}

		req, _ := http.NewRequest(http.MethodGet, "/jobs/other", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
