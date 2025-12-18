package sql

import (
	"context"
	"database/sql"
	"errors"

	"container-manager/internal/domain/entity"
	"container-manager/internal/infrastructure/database"
)

var _ database.JobDatabase = (*JobDatabase)(nil)

type JobDatabase struct {
	db *sql.DB
}

func NewJobDatabase(db *sql.DB) database.JobDatabase {
	return &JobDatabase{db: db}
}

func (db *JobDatabase) Create(ctx context.Context, job *entity.Job) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO jobs (id, type, status, payload, created_at, updated_at)	VALUES ($1, $2, $3, $4, $5, $6)",
		job.ID,
		job.Type,
		job.Status,
		job.Payload,
		job.CreatedAt,
		job.UpdatedAt,
	)
	return err
}

func (db *JobDatabase) GetByID(ctx context.Context, id string) (*entity.Job, error) {
	job := &entity.Job{}
	row := db.db.QueryRowContext(ctx, "SELECT id, type, status, payload, result, error, created_at, updated_at FROM jobs WHERE id = $1", id)
	err := row.Scan(
		&job.ID,
		&job.Type,
		&job.Status,
		&job.Payload,
		&job.Result,
		&job.Error,
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

func (db *JobDatabase) Update(ctx context.Context, job *entity.Job) error {
	_, err := db.db.ExecContext(ctx, "UPDATE jobs SET status = $2, result = $3, error = $4, updated_at = $5 WHERE id = $1",
		job.ID,
		job.Status,
		job.Result,
		job.Error,
		job.UpdatedAt,
	)
	return err
}
