package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rafaeljesus/srv-consumer/pkg"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
)

type (
	// UserStatusChanged is the message handler.
	UserStatusChanged struct {
		store pkg.UserStore
	}
)

// NewUserStatusChanged returns new UserStatusChanged struct.
func NewUserStatusChanged(s pkg.UserStore) *UserStatusChanged {
	return &UserStatusChanged{s}
}

// Handle is the user created message handler.
func (u *UserStatusChanged) Handle(ctx context.Context, m *message.Message) error {
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

	log.Print("user status successfully changed")
	return nil
}
