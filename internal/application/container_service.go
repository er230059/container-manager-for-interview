package application

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	containerruntime "container-manager/internal/infrastructure/container_runtime"

	"github.com/google/uuid"
)

type ContainerService struct {
	repo    repository.ContainerRepository
	jobRepo repository.JobRepository
}

func NewContainerService(repo repository.ContainerRepository, jobRepo repository.JobRepository) *ContainerService {
	return &ContainerService{
		repo:    repo,
		jobRepo: jobRepo,
	}
}

func (s *ContainerService) CreateContainer(ctx context.Context, userID int64, options containerruntime.ContainerCreateOptions) (string, error) {
	payload, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	job := &entity.Job{
		ID:        uuid.New().String(),
		Type:      "container_creation",
		Status:    entity.JobStatusPending,
		Payload:   payload,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.jobRepo.Create(context.Background(), job); err != nil {
		return "", err
	}

	go s.runCreateContainerJob(job, userID, options)

	return job.ID, nil
}

func (s *ContainerService) runCreateContainerJob(job *entity.Job, userID int64, options containerruntime.ContainerCreateOptions) {
	jobCtx := context.Background()

	job.Status = entity.JobStatusRunning
	job.UpdatedAt = time.Now()
	if err := s.jobRepo.Update(jobCtx, job); err != nil {
		log.Printf("failed to update job %s to running: %v", job.ID, err)
		return
	}

	container, err := s.repo.CreateContainer(jobCtx, userID, options)
	if err != nil {
		job.Status = entity.JobStatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		if updateErr := s.jobRepo.Update(jobCtx, job); updateErr != nil {
			log.Printf("failed to update job %s to failed: %v", job.ID, updateErr)
		}
		return
	}

	result, err := json.Marshal(map[string]string{"container_id": container.ID})
	if err != nil {
		job.Status = entity.JobStatusFailed
		job.Error = "failed to marshal result"
		job.UpdatedAt = time.Now()
	} else {
		job.Status = entity.JobStatusCompleted
		job.Result = result
		job.UpdatedAt = time.Now()
	}

	if updateErr := s.jobRepo.Update(jobCtx, job); updateErr != nil {
		log.Printf("failed to update job %s to completed/failed: %v", job.ID, updateErr)
	}
}

func (s *ContainerService) StartContainer(ctx context.Context, userID int64, id string) error {
	container, err := s.repo.GetContainer(ctx, id)
	if err != nil {
		return err
	}
	if container.UserID != userID {
		return errors.New("permission denied")
	}
	return s.repo.StartContainer(ctx, userID, id)
}

func (s *ContainerService) StopContainer(ctx context.Context, userID int64, id string) error {
	container, err := s.repo.GetContainer(ctx, id)
	if err != nil {
		return err
	}
	if container.UserID != userID {
		return errors.New("permission denied")
	}
	return s.repo.StopContainer(ctx, userID, id)
}

func (s *ContainerService) RemoveContainer(ctx context.Context, userID int64, id string) error {
	container, err := s.repo.GetContainer(ctx, id)
	if err != nil {
		return err
	}
	if container.UserID != userID {
		return errors.New("permission denied")
	}
	return s.repo.RemoveContainer(ctx, userID, id)
}
