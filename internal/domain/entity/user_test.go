package entity

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	id := int64(123)
	username := "testuser"
	plainPassword := "testpassword"

	t.Run("successful user creation", func(t *testing.T) {
		user, err := NewUser(id, username, plainPassword)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.ID != id {
			t.Errorf("expected ID %d, got %d", id, user.ID)
		}
		if user.Username != username {
			t.Errorf("expected username %s, got %s", username, user.Username)
		}

		// Verify password is hashed and not plain
		if user.Password == plainPassword {
			t.Error("password should be hashed")
		}

		// Verify hashed password can be validated
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
		if err != nil {
			t.Errorf("hashed password validation failed: %v", err)
		}
	})

	t.Run("NewUser with empty password", func(t *testing.T) {
		user, err := NewUser(id, username, "")
		if err == nil {
			t.Fatal("expected error for empty password, got nil")
		}
		if !errors.Is(err, ErrEmptyPassword) {
			t.Errorf("expected error %v, got %v", ErrEmptyPassword, err)
		}
		if user != nil {
			t.Errorf("expected nil user, got %v", user)
		}
	})

	// It's hard to directly test bcrypt.GenerateFromPassword error scenarios without
	// manipulating the environment (e.g., extremely high cost, but that makes test slow)
	// or mocking bcrypt itself. For now, assuming bcrypt works.
}

func TestUser_ValidatePassword(t *testing.T) {
	id := int64(456)
	username := "anotheruser"
	plainPassword := "securepassword"

	user, err := NewUser(id, username, plainPassword)
	if err != nil {
		t.Fatalf("failed to create user for validation test: %v", err)
	}

	t.Run("successful password validation", func(t *testing.T) {
		err := user.ValidatePassword(plainPassword)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("failed password validation (wrong password)", func(t *testing.T) {
		err := user.ValidatePassword("wrongpassword")
		if err == nil {
			t.Error("expected error for wrong password, got nil")
		}
		if err != bcrypt.ErrMismatchedHashAndPassword {
			t.Errorf("expected ErrMismatchedHashAndPassword, got %v", err)
		}
	})

	t.Run("failed password validation (empty password)", func(t *testing.T) {
		err := user.ValidatePassword("")
		if err == nil {
			t.Error("expected error for empty password, got nil")
		}
		if err != bcrypt.ErrMismatchedHashAndPassword {
			t.Errorf("expected ErrMismatchedHashAndPassword, got %v", err)
		}
	})
}
