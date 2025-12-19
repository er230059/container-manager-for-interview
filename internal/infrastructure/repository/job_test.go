package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"container-manager/internal/domain/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestJobRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewJobRepository(db)
	ctx := context.Background()

	now := time.Now()
	job := &entity.Job{
		ID:        "job-1",
		Type:      "test-job",
		Status:    "pending",
		Payload:   json.RawMessage([]byte("{}")),
		UserID:    123,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO jobs").
			WithArgs(job.ID, job.Type, job.Status, job.Payload, job.UserID, job.CreatedAt, job.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, job)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO jobs").
			WithArgs(job.ID, job.Type, job.Status, job.Payload, job.UserID, job.CreatedAt, job.UpdatedAt).
			WillReturnError(errors.New("db error"))

		err := repo.Create(ctx, job)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJobRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewJobRepository(db)
	ctx := context.Background()

	now := time.Now()
	job := &entity.Job{
		ID:        "job-1",
		Type:      "test-job",
		Status:    "pending",
		Payload:   json.RawMessage([]byte("{}")),
		Result:    json.RawMessage([]byte("null")),
		Error:     "",
		UserID:    123,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
	
rows := sqlmock.NewRows([]string{"id", "type", "status", "payload", "result", "error", "user_id", "created_at", "updated_at"}).
			AddRow(job.ID, job.Type, job.Status, job.Payload, job.Result, job.Error, job.UserID, job.CreatedAt, job.UpdatedAt)

		mock.ExpectQuery("SELECT id, type, status, payload, result, error, user_id, created_at, updated_at FROM jobs WHERE id = \\$1").
			WithArgs(job.ID).
			WillReturnRows(rows)

		result, err := repo.GetByID(ctx, job.ID)
		assert.NoError(t, err)
		assert.Equal(t, job.ID, result.ID)
		assert.Equal(t, job.Type, result.Type)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, type, status, payload, result, error, user_id, created_at, updated_at FROM jobs WHERE id = \\$1").
			WithArgs("non-existent").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetByID(ctx, "non-existent")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, type, status, payload, result, error, user_id, created_at, updated_at FROM jobs WHERE id = \\$1").
			WithArgs(job.ID).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetByID(ctx, job.ID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJobRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewJobRepository(db)
	ctx := context.Background()

	now := time.Now()
	job := &entity.Job{
		ID:        "job-1",
		Status:    "completed",
		Result:    json.RawMessage([]byte(`"success"`)),
		Error:     "",
		UpdatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE jobs SET status = \\$2, result = \\$3, error = \\$4, updated_at = \\$5 WHERE id = \\$1").
			WithArgs(job.ID, job.Status, job.Result, job.Error, job.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(ctx, job)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec("UPDATE jobs SET status = \\$2, result = \\$3, error = \\$4, updated_at = \\$5 WHERE id = \\$1").
			WithArgs(job.ID, job.Status, job.Result, job.Error, job.UpdatedAt).
			WillReturnError(errors.New("db error"))

		err := repo.Update(ctx, job)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
