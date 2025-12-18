package application

import (
	"context"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
)

type JobService interface {
	GetJob(ctx context.Context, id string) (*entity.Job, error)
}

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) JobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}

func (s *jobService) GetJob(ctx context.Context, id string) (*entity.Job, error) {
	return s.jobRepo.GetByID(ctx, id)
}
