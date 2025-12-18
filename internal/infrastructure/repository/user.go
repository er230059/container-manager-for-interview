package repository

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/infrastructure"
	"context"
	"database/sql"
	"errors"
)

var _ infrastructure.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(
	db *sql.DB,
) infrastructure.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Password)
	return err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"
	row := r.db.QueryRowContext(ctx, query, username)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
