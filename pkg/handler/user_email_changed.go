package handler

import (
	"encoding/json"
	"log"

	"github.com/rafaeljesus/srv-consumer/pkg"
	"github.com/rafaeljesus/srv-consumer/pkg/errors"
	"github.com/rafaeljesus/srv-consumer/pkg/message"
)

type (
	UserEmailChanged struct {
		store pkg.UserStore
	}
)

func NewUserEmailChanged(s pkg.UserStore) *UserEmailChanged {
	return &UserEmailChanged{s}
}

func (u *UserEmailChanged) Handle(m *message.Message) error {
	defer m.Ack(true)

	user := new(pkg.User)
	if err := json.Unmarshal(m.Body, user); err != nil {
		log.Printf("failed to unmarshal message body: %v", err)
		return err
	}

	if err := u.store.Save(user); err != nil {
		if err == errors.ErrNotFound {
			log.Print("user not found")
		} else {
			log.Printf("failed to save user to store: %v", err)
		}
		return err
	}

	log.Print("user email successfully changed")
	return nil
}
