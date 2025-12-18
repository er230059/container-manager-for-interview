package repository

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	database "container-manager/internal/infrastructure/database/sql"
	"context"
)

var _ repository.UserRepository = (*userRepository)(nil)

type userRepository struct {
	userDatabase *database.UserDatabase
}

func NewUserRepository(
	userDatabase *database.UserDatabase,
) repository.UserRepository {
	return &userRepository{
		userDatabase: userDatabase,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.userDatabase.Create(ctx, user)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return r.userDatabase.FindByUsername(ctx, username)

}
