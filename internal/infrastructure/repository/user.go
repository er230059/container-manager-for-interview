package repository

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	"container-manager/internal/infrastructure/database"
	"context"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db database.UserDatabase
}

func NewUserRepository(db database.UserDatabase) repository.UserRepository {
	return &UserRepository{db: db}
}

// Create saves a new user to the database.
func (r *UserRepository) Create(user *entity.User) error {
	// Here we could add a context with timeout.
	ctx := context.Background()
	return r.db.CreateUser(ctx, user.ID, user.Username, user.Password)
}

// FindByUsername retrieves a user by their username.
func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	// Here we could add a context with timeout.
	ctx := context.Background()
	return r.db.FindByUsername(ctx, username)
}
