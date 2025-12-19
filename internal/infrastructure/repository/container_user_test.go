package repository

import (
	"context"
	"database/sql"
	"testing"

	internalErrors "container-manager/internal/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestContainerUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewContainerUserRepository(db)
	ctx := context.Background()

	containerID := "container-1"
	userID := int64(123)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO container_user").
			WithArgs(containerID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, containerID, userID)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO container_user").
			WithArgs(containerID, userID).
			WillReturnError(sql.ErrConnDone)

		err := repo.Create(ctx, containerID, userID)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContainerUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewContainerUserRepository(db)
	ctx := context.Background()

	containerID := "container-1"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM container_user WHERE container_id = \\$1").
			WithArgs(containerID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, containerID)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM container_user WHERE container_id = \\$1").
			WithArgs(containerID).
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(ctx, containerID)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContainerUserRepository_GetUserIDByContainerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewContainerUserRepository(db)
	ctx := context.Background()

	containerID := "container-1"
	userID := int64(123)

	t.Run("success", func(t *testing.T) {
	
rows := sqlmock.NewRows([]string{"user_id"}).AddRow(userID)
		mock.ExpectQuery("SELECT user_id FROM container_user WHERE container_id = \\$1").
			WithArgs(containerID).
			WillReturnRows(rows)

		result, err := repo.GetUserIDByContainerID(ctx, containerID)
		assert.NoError(t, err)
		assert.Equal(t, userID, result)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT user_id FROM container_user WHERE container_id = \\$1").
			WithArgs(containerID).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetUserIDByContainerID(ctx, containerID)
		assert.Error(t, err)
		assert.Equal(t, internalErrors.ContainerNotFound, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT user_id FROM container_user WHERE container_id = \\$1").
			WithArgs(containerID).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.GetUserIDByContainerID(ctx, containerID)
		assert.Error(t, err)
		assert.NotEqual(t, internalErrors.ContainerNotFound, err)
		assert.Equal(t, int64(0), result)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContainerUserRepository_GetContainerIDsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewContainerUserRepository(db)
	ctx := context.Background()

	userID := int64(123)
	containerID1 := "container-1"
	containerID2 := "container-2"

	t.Run("success", func(t *testing.T) {
	
rows := sqlmock.NewRows([]string{"container_id"}).
			AddRow(containerID1).
			AddRow(containerID2)

		mock.ExpectQuery("SELECT container_id FROM container_user WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetContainerIDsByUserID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, []string{containerID1, containerID2}, result)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT container_id FROM container_user WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.GetContainerIDsByUserID(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
