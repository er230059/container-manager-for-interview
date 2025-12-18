package repository

import (
	"context"

	"container-manager/internal/domain/entity"
)

type JobRepository interface {
	Create(ctx context.Context, job *entity.Job) error
	GetByID(ctx context.Context, id string) (*entity.Job, error)
	Update(ctx context.Context, job *entity.Job) error
}
