package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/rafaeljesus/srv-consumer"
	"github.com/rafaeljesus/srv-consumer/mock"
	"github.com/rafaeljesus/srv-consumer/platform/message"
)

var (
	errAcker = errors.New("acker error")
)

func TestUserCreated(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *mock.UserStore, *mock.Acknowledger)
	}{
		{
			"handle user created",
			testHandleUserCreated,
		},
		{
			"fail to unmarshal body",
			testFailToUnmarshalBody,
		},
		{
			"fail to ack when unmarshal body error",
			testFailAckWhenUnmarshalBodyError,
		},
		{
			"handle conflict error",
			testHandleConflictError,
		},
		{
			"handle unexpected error",
			testHandleUnexpectedError,
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
	store.AddFunc = func(user *srv.User) error {
		if user.Email != "foo@mail.com" {
			t.Fatal("unexpected email")
		}
		if user.Username != "foo" {
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

func testFailToUnmarshalBody(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *srv.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return nil }
	body := []byte(``)

	msg := message.New(acker, body)
	h := NewUserCreated(store)
	err := h.Handle(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected to return err: %v", err)
	}
	if store.AddInvoked {
		t.Fatal("expected store.Add() to not be invoked")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be invoked")
	}
}

func testFailAckWhenUnmarshalBodyError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *srv.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return errAcker }
	body := []byte(``)

	msg := message.New(acker, body)
	h := NewUserCreated(store)
	err := h.Handle(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected to return err: %v", err)
	}
	if store.AddInvoked {
		t.Fatal("expected store.Add() to not be invoked")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be invoked")
	}
}

func testHandleConflictError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *srv.User) error { return srv.ErrConflict }
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
	if err != srv.ErrConflict {
		t.Fatalf("expected to return err: %v", err)
	}
	if !store.AddInvoked {
		t.Fatal("expected store.Add() to not be invoked")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be invoked")
	}
}

func testHandleUnexpectedError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *srv.User) error { return errors.New("unexpected error") }
	acker.NackFunc = func(multiple, requeue bool) error {
		if multiple {
			t.Fatal("unexpected multiple")
		}
		if !requeue {
			t.Fatal("unexpected requeue")
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
	if err == nil {
		t.Fatalf("expected to return err: %v", err)
	}
	if !store.AddInvoked {
		t.Fatal("expected store.Add() to not be invoked")
	}
	if !acker.NackInvoked {
		t.Fatal("expected message.Nack() to be invoked")
	}
}
