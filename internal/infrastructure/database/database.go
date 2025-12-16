package database

import (
	"context"
)

// UserDatabase defines the interface for user data operations,
// abstracting away the specific database technology.
type UserDatabase interface {
	CreateUser(ctx context.Context, id int64, username, password string) error
}
