package repository

import (
	"container-manager/internal/domain"
	"container-manager/internal/infrastructure/database"
	"context"
)

// UserRepository is an implementation of the domain.UserRepository interface.
type UserRepository struct {
	db database.UserDatabase
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db database.UserDatabase) domain.UserRepository {
	return &UserRepository{db: db}
}

// Create saves a new user to the database.
func (r *UserRepository) Create(user *domain.User) error {
	// Here we could add a context with timeout.
	ctx := context.Background()
	return r.db.CreateUser(ctx, user.ID, user.Username, user.Password)
}
