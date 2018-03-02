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
	// UserCreated is the message handler.
	UserCreated struct {
		store pkg.UserStore
	}
)

// NewUserCreated returns new UserCreated struct.
func NewUserCreated(s pkg.UserStore) *UserCreated {
	return &UserCreated{s}
}

// Handle is the user created message handler.
func (u *UserCreated) Handle(ctx context.Context, m *message.Message) error {
	user := new(pkg.User)
	if err := json.Unmarshal(m.Body, user); err != nil {
		log.Printf("failed to unmarshal message body: %v", err)
		if err := m.Ack(false); err != nil {
			log.Printf("failed to ack message: %v", err)
		}
		return err
	}

	err := u.store.Add(user)
	switch err {
	case nil:
		log.Print("user successfully added")
		if err := m.Ack(false); err != nil {
			return fmt.Errorf("failed to ack message: %v", err)
		}
		return nil
	case pkg.ErrConflict:
		log.Print("user already exists")
		if err := m.Ack(false); err != nil {
			log.Printf("failed to nack message: %v", err)
		}
		return err
	default:
		log.Printf("failed to add user to store: %v", err)
		if err := m.Nack(false, true); err != nil {
			log.Printf("failed to reject message: %v", err)
		}
		return err
	}
}
