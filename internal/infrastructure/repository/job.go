package repository

import (
	"context"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	"container-manager/internal/infrastructure/database"
)

type jobRepository struct {
	db database.JobDatabase
}

func NewJobRepository(db database.JobDatabase) repository.JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *entity.Job) error {
	return r.db.Create(ctx, job)
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*entity.Job, error) {
	return r.db.GetByID(ctx, id)
}

func (r *jobRepository) Update(ctx context.Context, job *entity.Job) error {
	return r.db.Update(ctx, job)
}
