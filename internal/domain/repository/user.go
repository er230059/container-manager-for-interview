package repository

import "container-manager/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	// FindByUsername(username string) (*User, error) // Example for future use
}
