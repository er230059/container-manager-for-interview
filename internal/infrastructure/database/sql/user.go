package sql

import (
	"container-manager/internal/infrastructure/database"
	"context"
	"database/sql"
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
