package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

// Handle is the user status changed message handler.
func (u *UserStatusChanged) Handle(ctx context.Context, m *message.Message) error {
	user := new(pkg.User)
	if err := json.Unmarshal(m.Body, user); err != nil {
		log.Printf("failed to unmarshal message body: %v", err)
		if err := m.Ack(false); err != nil {
			log.Printf("failed to ack message: %v", err)
			return err
		}
		return err
	}

	err := u.store.Save(user)

	switch err {
	case nil:
		log.Print("user status successfully changed")
		if err := m.Ack(false); err != nil {
			return fmt.Errorf("failed to ack message: %v", err)
		}
		return nil
	case pkg.ErrNotFound:
		log.Print("user not found")
		if err := m.Ack(false); err != nil {
			return fmt.Errorf("failed to ack message: %v", err)
		}
		return err
	default:
		log.Printf("failed to save user to store: %v", err)
		if err := m.Nack(false, true); err != nil {
			log.Printf("failed to reject message: %v", err)
		}
		return err
	}
}
