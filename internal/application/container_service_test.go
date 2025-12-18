package application

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"container-manager/internal/application/mocks"
	"container-manager/internal/domain/entity"
	containerruntime "container-manager/internal/infrastructure/container_runtime"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestContainerService_CreateContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	options := containerruntime.ContainerCreateOptions{
		Image: "test-image",
	}

	var wg sync.WaitGroup
	wg.Add(1)

	mockJobRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, job *entity.Job) error {
		assert.Equal(t, "container_creation", job.Type)
		assert.Equal(t, entity.JobStatusPending, job.Status)
		assert.Equal(t, userID, job.UserID)
		return nil
	})

	// Expectations for the goroutine
	gomock.InOrder(
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
		mockContainerRepo.EXPECT().CreateContainer(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Container{ID: "some-id"}, nil),
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ *entity.Job) error {
			wg.Done()
			return nil
		}),
	)

	jobID, err := service.CreateContainer(ctx, userID, options)
	assert.NoError(t, err)
	assert.NotEmpty(t, jobID)

	wg.Wait()
}

func TestContainerService_runCreateContainerJob_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	userID := int64(1)
	options := containerruntime.ContainerCreateOptions{
		Image: "test-image",
	}
	payload, _ := json.Marshal(options)

	job := &entity.Job{
		ID:      "job-123",
		Type:    "container_creation",
		Status:  entity.JobStatusPending,
		Payload: payload,
		UserID:  userID,
	}

	container := &entity.Container{
		ID: "container-123",
	}

	gomock.InOrder(
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, job *entity.Job) error {
			assert.Equal(t, entity.JobStatusRunning, job.Status)
			return nil
		}),
		mockContainerRepo.EXPECT().CreateContainer(gomock.Any(), userID, options).Return(container, nil),
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, job *entity.Job) error {
			assert.Equal(t, entity.JobStatusCompleted, job.Status)
			var result map[string]string
			err := json.Unmarshal(job.Result, &result)
			assert.NoError(t, err)
			assert.Equal(t, container.ID, result["container_id"])
			return nil
		}),
	)

	service.runCreateContainerJob(job, userID, options)
}

func TestContainerService_runCreateContainerJob_CreateContainerFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	userID := int64(1)
	options := containerruntime.ContainerCreateOptions{
		Image: "test-image",
	}
	payload, _ := json.Marshal(options)

	job := &entity.Job{
		ID:      "job-123",
		Type:    "container_creation",
		Status:  entity.JobStatusPending,
		Payload: payload,
		UserID:  userID,
	}

	createErr := errors.New("create container error")

	gomock.InOrder(
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, job *entity.Job) error {
			assert.Equal(t, entity.JobStatusRunning, job.Status)
			return nil
		}),
		mockContainerRepo.EXPECT().CreateContainer(gomock.Any(), userID, options).Return(nil, createErr),
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, job *entity.Job) error {
			assert.Equal(t, entity.JobStatusFailed, job.Status)
			assert.Equal(t, createErr.Error(), job.Error)
			return nil
		}),
	)

	service.runCreateContainerJob(job, userID, options)
}

func TestContainerService_StartContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	containerID := "container-123"

	container := &entity.Container{
		ID:     containerID,
		UserID: userID,
	}

	mockContainerRepo.EXPECT().GetContainer(ctx, containerID).Return(container, nil)
	mockContainerRepo.EXPECT().StartContainer(ctx, userID, containerID).Return(nil)

	err := service.StartContainer(ctx, userID, containerID)
	assert.NoError(t, err)
}

func TestContainerService_StartContainer_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	otherUserID := int64(2)
	containerID := "container-123"

	container := &entity.Container{
		ID:     containerID,
		UserID: otherUserID,
	}

	mockContainerRepo.EXPECT().GetContainer(ctx, containerID).Return(container, nil)

	err := service.StartContainer(ctx, userID, containerID)
	assert.EqualError(t, err, "permission denied")
}

func TestContainerService_StopContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	containerID := "container-123"

	container := &entity.Container{
		ID:     containerID,
		UserID: userID,
	}

	mockContainerRepo.EXPECT().GetContainer(ctx, containerID).Return(container, nil)
	mockContainerRepo.EXPECT().StopContainer(ctx, userID, containerID).Return(nil)

	err := service.StopContainer(ctx, userID, containerID)
	assert.NoError(t, err)
}

func TestContainerService_StopContainer_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	otherUserID := int64(2)
	containerID := "container-123"

	container := &entity.Container{
		ID:     containerID,
		UserID: otherUserID,
	}

	mockContainerRepo.EXPECT().GetContainer(ctx, containerID).Return(container, nil)

	err := service.StopContainer(ctx, userID, containerID)
	assert.EqualError(t, err, "permission denied")
}

func TestContainerService_RemoveContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	containerID := "container-123"

	container := &entity.Container{
		ID:     containerID,
		UserID: userID,
	}

	mockContainerRepo.EXPECT().GetContainer(ctx, containerID).Return(container, nil)
	mockContainerRepo.EXPECT().RemoveContainer(ctx, userID, containerID).Return(nil)

	err := service.RemoveContainer(ctx, userID, containerID)
	assert.NoError(t, err)
}

func TestContainerService_RemoveContainer_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainerRepo := mocks.NewMockContainerRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockContainerRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	otherUserID := int64(2)
	containerID := "container-123"

	container := &entity.Container{
		ID:     containerID,
		UserID: otherUserID,
	}

	mockContainerRepo.EXPECT().GetContainer(ctx, containerID).Return(container, nil)

	err := service.RemoveContainer(ctx, userID, containerID)
	assert.EqualError(t, err, "permission denied")
}
