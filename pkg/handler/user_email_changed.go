package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rafaeljesus/srv-consumer/pkg"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
)

type (
	// UserEmailChanged is the message handler.
	UserEmailChanged struct {
		store pkg.UserStore
	}
)

// NewUserEmailChanged returns new UserEmailChanged struct.
func NewUserEmailChanged(s pkg.UserStore) *UserEmailChanged {
	return &UserEmailChanged{s}
}

// Handle is the user email changed message handler.
func (u *UserEmailChanged) Handle(ctx context.Context, m *message.Message) error {
	defer m.Ack(true)

	user := new(pkg.User)
	if err := json.Unmarshal(m.Body, user); err != nil {
		log.Printf("failed to unmarshal message body: %v", err)
		return err
	}

	if err := u.store.Save(user); err != nil {
		if err == pkg.ErrNotFound {
			log.Print("user not found")
		} else {
			log.Printf("failed to save user to store: %v", err)
		}
		return err
	}

	log.Print("user email successfully changed")
	return nil
}
