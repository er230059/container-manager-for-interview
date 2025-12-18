package repository

import (
	"context"
	"database/sql"
	"errors"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/infrastructure"
)

type jobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) infrastructure.JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *entity.Job) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO jobs (id, type, status, payload, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		job.ID,
		job.Type,
		job.Status,
		job.Payload,
		job.UserID,
		job.CreatedAt,
		job.UpdatedAt,
	)
	return err
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*entity.Job, error) {
	job := &entity.Job{}
	row := r.db.QueryRowContext(ctx, "SELECT id, type, status, payload, result, error, user_id, created_at, updated_at FROM jobs WHERE id = $1", id)
	err := row.Scan(
		&job.ID,
		&job.Type,
		&job.Status,
		&job.Payload,
		&job.Result,
		&job.Error,
		&job.UserID,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return job, nil
}

func (r *jobRepository) Update(ctx context.Context, job *entity.Job) error {
	_, err := r.db.ExecContext(ctx, "UPDATE jobs SET status = $2, result = $3, error = $4, updated_at = $5 WHERE id = $1",
		job.ID,
		job.Status,
		job.Result,
		job.Error,
		job.UpdatedAt,
	)
	return err
}
