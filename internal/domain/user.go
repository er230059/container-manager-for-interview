package domain

// User represents a user in the system.
type User struct {
	ID       int64
	Username string
	Password string // In a real app, this should be a hashed value.
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(user *User) error
	// FindByUsername(username string) (*User, error) // Example for future use
}