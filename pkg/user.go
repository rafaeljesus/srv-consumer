package pkg

import "errors"

var (
	// ErrConflict is the conflict error.
	ErrConflict = errors.New("conflict error")
	// ErrNotFound is the not found error.
	ErrNotFound = errors.New("not found")
)

type (
	// User is the model struct which represents a user.
	User struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Status   string `json:"status"`
	}

	// UserStore contains methods for managing users in a storage.
	UserStore interface {
		// Add a new user to the store.
		Add(user *User) error
		// Save a user to the store.
		Save(user *User) error
	}
)
