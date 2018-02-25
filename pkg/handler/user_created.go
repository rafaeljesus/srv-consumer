package handler

import (
	"encoding/json"
	"log"

	"github.com/rafaeljesus/srv-consumer/pkg"
	"github.com/rafaeljesus/srv-consumer/pkg/errors"
	"github.com/rafaeljesus/srv-consumer/pkg/message"
)

type (
	UserCreated struct {
		store pkg.UserStore
	}
)

func NewUserCreated(s pkg.UserStore) *UserCreated {
	return &UserCreated{s}
}

func (u *UserCreated) Handle(m *message.Message) error {
	defer m.Ack(false)

	user := new(pkg.User)
	if err := json.Unmarshal(m.Body, user); err != nil {
		log.Printf("failed to unmarshal message body: %v", err)
		return err
	}

	err := u.store.Add(user)
	switch err {
	case nil:
		log.Print("user successfully added")
		return nil
	case errors.ErrConflict:
		log.Print("user already exists")
		return err
	default:
		log.Printf("failed to add user to store: %v", err)
		return err
	}
}
