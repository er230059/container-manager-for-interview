package application

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"container-manager/internal/application/mocks"
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestContainerService_CreateContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)

	ctx := context.Background()
	userID := int64(1)
	options := infrastructure.ContainerCreateOptions{
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
		mockRuntime.EXPECT().Create(gomock.Any(), options).Return("container-123", nil),
		mockContainerUserRepo.EXPECT().Create(gomock.Any(), "container-123", userID).Return(nil),
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ *entity.Job) error {
			wg.Done()
			return nil
		}),
	)

	jobID, err := service.CreateContainer(ctx, userID, options)
	assert.NoError(t, err)
	assert.NotEmpty(t, jobID)

	wg.Wait() // Wait for the goroutine to finish
}

func TestContainerService_runCreateContainerJob_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)

	userID := int64(1)
	options := infrastructure.ContainerCreateOptions{
		Image: "test-image",
	}
	payload, _ := json.Marshal(options)
	job := &entity.Job{
		ID:        uuid.NewString(),
		Type:      "container_creation",
		Status:    entity.JobStatusPending,
		Payload:   payload,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	containerID := "container-123"

	gomock.InOrder(
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedJob *entity.Job) error {
			assert.Equal(t, entity.JobStatusRunning, updatedJob.Status)
			return nil
		}),
		mockRuntime.EXPECT().Create(gomock.Any(), options).Return(containerID, nil),
		mockContainerUserRepo.EXPECT().Create(gomock.Any(), containerID, userID).Return(nil),
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedJob *entity.Job) error {
			assert.Equal(t, entity.JobStatusCompleted, updatedJob.Status)
			var result map[string]string
			err := json.Unmarshal(updatedJob.Result, &result)
			assert.NoError(t, err)
			assert.Equal(t, containerID, result["container_id"])
			return nil
		}),
	)

	service.runCreateContainerJob(job, userID, options)
}

func TestContainerService_runCreateContainerJob_RuntimeCreateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)

	userID := int64(1)
	options := infrastructure.ContainerCreateOptions{
		Image: "test-image",
	}
	payload, _ := json.Marshal(options)
	job := &entity.Job{
		ID:      uuid.NewString(),
		UserID:  userID,
		Payload: payload,
	}
	createErr := errors.New("runtime create error")

	gomock.InOrder(
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
		mockRuntime.EXPECT().Create(gomock.Any(), options).Return("", createErr),
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedJob *entity.Job) error {
			assert.Equal(t, entity.JobStatusFailed, updatedJob.Status)
			assert.Equal(t, createErr.Error(), updatedJob.Error)
			return nil
		}),
	)

	service.runCreateContainerJob(job, userID, options)
}

func TestContainerService_runCreateContainerJob_UserRepoCreateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)
	mockJobRepo := mocks.NewMockJobRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, mockJobRepo)

	userID := int64(1)
	options := infrastructure.ContainerCreateOptions{
		Image: "test-image",
	}
	payload, _ := json.Marshal(options)
	job := &entity.Job{
		ID:      uuid.NewString(),
		UserID:  userID,
		Payload: payload,
	}
	containerID := "container-123"
	repoErr := errors.New("user repo create error")

	gomock.InOrder(
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
		mockRuntime.EXPECT().Create(gomock.Any(), options).Return(containerID, nil),
		mockContainerUserRepo.EXPECT().Create(gomock.Any(), containerID, userID).Return(repoErr),
		mockRuntime.EXPECT().Remove(gomock.Any(), containerID).Return(nil), // Rollback
		mockJobRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, updatedJob *entity.Job) error {
			assert.Equal(t, entity.JobStatusFailed, updatedJob.Status)
			assert.Equal(t, repoErr.Error(), updatedJob.Error)
			return nil
		}),
	)

	service.runCreateContainerJob(job, userID, options)
}

func TestContainerService_StartContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	containerID := "container-123"

	mockContainerUserRepo.EXPECT().GetUserIDByContainerID(ctx, containerID).Return(userID, nil)
	mockRuntime.EXPECT().Start(ctx, containerID).Return(nil)

	err := service.StartContainer(ctx, userID, containerID)
	assert.NoError(t, err)
}

func TestContainerService_StartContainer_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	otherUserID := int64(2)
	containerID := "container-123"

	mockContainerUserRepo.EXPECT().GetUserIDByContainerID(ctx, containerID).Return(otherUserID, nil)

	err := service.StartContainer(ctx, userID, containerID)
	assert.EqualError(t, err, "permission denied")
}

func TestContainerService_StopContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	containerID := "container-123"

	mockContainerUserRepo.EXPECT().GetUserIDByContainerID(ctx, containerID).Return(userID, nil)
	mockRuntime.EXPECT().Stop(ctx, containerID).Return(nil)

	err := service.StopContainer(ctx, userID, containerID)
	assert.NoError(t, err)
}

func TestContainerService_StopContainer_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	otherUserID := int64(2)
	containerID := "container-123"

	mockContainerUserRepo.EXPECT().GetUserIDByContainerID(ctx, containerID).Return(otherUserID, nil)

	err := service.StopContainer(ctx, userID, containerID)
	assert.EqualError(t, err, "permission denied")
}

func TestContainerService_RemoveContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	containerID := "container-123"

	mockContainerUserRepo.EXPECT().GetUserIDByContainerID(ctx, containerID).Return(userID, nil)
	mockContainerUserRepo.EXPECT().Delete(ctx, containerID).Return(nil)
	mockRuntime.EXPECT().Remove(ctx, containerID).Return(nil)

	err := service.RemoveContainer(ctx, userID, containerID)
	assert.NoError(t, err)
}

func TestContainerService_RemoveContainer_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	otherUserID := int64(2)
	containerID := "container-123"

	mockContainerUserRepo.EXPECT().GetUserIDByContainerID(ctx, containerID).Return(otherUserID, nil)

	err := service.RemoveContainer(ctx, userID, containerID)
	assert.EqualError(t, err, "permission denied")
}

func TestContainerService_ListContainers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	containerID1 := "container-1"
	containerID2 := "container-2"

	expectedContainer1 := &entity.Container{ID: containerID1, Image: "test-image-1"}
	expectedContainer2 := &entity.Container{ID: containerID2, Image: "test-image-2"}

	mockContainerUserRepo.EXPECT().GetContainerIDsByUserID(ctx, userID).Return([]string{containerID1, containerID2}, nil)
	mockRuntime.EXPECT().Inspect(ctx, containerID1).Return(expectedContainer1, nil)
	mockRuntime.EXPECT().Inspect(ctx, containerID2).Return(expectedContainer2, nil)

	containers, err := service.ListContainers(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, containers, 2)
	assert.Equal(t, expectedContainer1, containers[0])
	assert.Equal(t, expectedContainer2, containers[1])
}

func TestContainerService_ListContainers_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	repoErr := errors.New("repo error")

	mockContainerUserRepo.EXPECT().GetContainerIDsByUserID(ctx, userID).Return(nil, repoErr)

	containers, err := service.ListContainers(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, repoErr, err)
	assert.Nil(t, containers)
}

func TestContainerService_ListContainers_InspectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntime(ctrl)
	mockContainerUserRepo := mocks.NewMockContainerUserRepository(ctrl)

	service := NewContainerService(mockRuntime, mockContainerUserRepo, nil)

	ctx := context.Background()
	userID := int64(1)
	containerID1 := "container-1"
	containerID2 := "container-2"

	expectedContainer1 := &entity.Container{ID: containerID1, Image: "test-image-1"}
	inspectErr := errors.New("inspect error")

	mockContainerUserRepo.EXPECT().GetContainerIDsByUserID(ctx, userID).Return([]string{containerID1, containerID2}, nil)
	mockRuntime.EXPECT().Inspect(ctx, containerID1).Return(expectedContainer1, nil)
	mockRuntime.EXPECT().Inspect(ctx, containerID2).Return(nil, inspectErr)

	containers, err := service.ListContainers(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, containers, 1)
	assert.Equal(t, expectedContainer1, containers[0])
}
