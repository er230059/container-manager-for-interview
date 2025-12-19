package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"container-manager/internal/domain/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Username, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Username, user.Password).
			WillReturnError(errors.New("db error"))

		err := repo.Create(ctx, user)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}

	t.Run("success", func(t *testing.T) {
	
rows := sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(user.ID, user.Username, user.Password)

		mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = \\$1").
			WithArgs(user.Username).
			WillReturnRows(rows)

		result, err := repo.FindByUsername(ctx, user.Username)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Username, result.Username)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = \\$1").
			WithArgs("non-existent").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByUsername(ctx, "non-existent")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = \\$1").
			WithArgs(user.Username).
			WillReturnError(errors.New("db error"))

		result, err := repo.FindByUsername(ctx, user.Username)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
