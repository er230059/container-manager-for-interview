package sql

import (
	"container-manager/internal/infrastructure/database"
	"context"
	"database/sql"
)

// UserDatabase is a PostgreSQL implementation of the UserDatabase interface.
type UserDatabase struct {
	db *sql.DB
}

// Compile-time check to ensure UserDatabase implements UserDatabase.
var _ database.UserDatabase = (*UserDatabase)(nil)

// NewUserDatabase creates a new Postgres-backed UserDatabase.
func NewUserDatabase(db *sql.DB) database.UserDatabase {
	return &UserDatabase{db: db}
}

// CreateUser inserts a new user into the database.
func (d *UserDatabase) CreateUser(ctx context.Context, id int64, username, password string) error {
	query := "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)"
	_, err := d.db.ExecContext(ctx, query, id, username, password)
	return err
}
