package handler

import (
	"context"
	"testing"

	"github.com/rafaeljesus/srv-consumer/pkg"
	"github.com/rafaeljesus/srv-consumer/pkg/mock"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
)

func TestUserCreated(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *mock.UserStore, *mock.Acknowledger)
	}{
		{
			"handler user created",
			testHandleUserCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			s := new(mock.UserStore)
			a := new(mock.Acknowledger)
			test.function(t, s, a)
		})
	}
}

func testHandleUserCreated(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *pkg.User) error {
		if user.Email == "" {
			t.Fatal("unexpected email")
		}
		if user.Username == "" {
			t.Fatal("unexpected username")
		}
		return nil
	}
	acker.AckFunc = func(multiple bool) error {
		if multiple {
			t.Fatal("unexpected multiple")
		}
		return nil
	}
	body := []byte(`{
		"email": "foo@mail.com",
		"username": "foo",
		"status": "new"
	}`)

	msg := message.New(acker, body)
	h := NewUserCreated(store)
	err := h.Handle(context.Background(), msg)
	if err != nil {
		t.Fatalf("expected to handle user created %v", err)
	}
	if !store.AddInvoked {
		t.Fatal("expected store.Add() to be invoked")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be invoked")
	}
}
