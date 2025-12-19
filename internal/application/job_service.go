package application

import (
	"container-manager/internal/errors"
	"context"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/infrastructure"
)

type JobService interface {
	GetJob(ctx context.Context, userID int64, id string) (*entity.Job, error)
}

type jobService struct {
	jobRepo infrastructure.JobRepository
}

func NewJobService(jobRepo infrastructure.JobRepository) JobService {
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
		return nil, errors.PermissionDenied
	}
	return job, nil
}
