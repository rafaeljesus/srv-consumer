package pkg

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
		Add(user *User) error
		Save(user *User) error
	}
)
