package sql

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/infrastructure/database"
	"context"
	"database/sql"
	"errors"
)

var _ database.UserDatabase = (*UserDatabase)(nil)

type UserDatabase struct {
	db *sql.DB
}

func NewUserDatabase(db *sql.DB) database.UserDatabase {
	return &UserDatabase{db: db}
}

func (d *UserDatabase) CreateUser(ctx context.Context, id int64, username, password string) error {
	query := "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)"
	_, err := d.db.ExecContext(ctx, query, id, username, password)
	return err
}

func (d *UserDatabase) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"
	row := d.db.QueryRowContext(ctx, query, username)
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
