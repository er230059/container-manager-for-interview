package application

import (
	"context"
	"errors"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
)

type JobService interface {
	GetJob(ctx context.Context, userID int64, id string) (*entity.Job, error)
}

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) JobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}

func (s *jobService) GetJob(ctx context.Context, userID int64, id string) (*entity.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if job.UserID != userID {
		return nil, errors.New("permission denied")
	}
	return job, nil
}
