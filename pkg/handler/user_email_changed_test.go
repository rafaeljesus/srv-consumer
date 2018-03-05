package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/rafaeljesus/srv-consumer/pkg"
	"github.com/rafaeljesus/srv-consumer/pkg/mock"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
)

func TestUserEmailChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, *mock.UserStore, *mock.Acknowledger)
	}{
		{
			scenario: "When valid payload is supplied, then should successfully save user",
			function: testShouldSuccessfullyChangeUserEmail,
		},
		{
			scenario: "When invalid payload is supplied, then should fail to unmarshal body",
			function: testShouldFailToUnmarshalBody,
		},
		{
			scenario: "When Not found user is supplied, then should fail to save",
			function: testHandleNotFoundError,
		},

		{
			scenario: "When unexpected error occurs, should be handled properly",
			function: testHandleUnexpectedSaveError,
		},
		{
			scenario: "When unable to Ack message, error should be handled properly",
			function: testEmailChangeHandlerShouldFailToAck,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			store := new(mock.UserStore)
			acker := new(mock.Acknowledger)
			test.function(t, store, acker)

		})
	}
}

func testShouldSuccessfullyChangeUserEmail(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *pkg.User) error {
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
        "status": "active"
	}`)

	msg := message.New(acker, body)
	handler := NewUserEmailChanged(store)
	err := handler.Handle(context.Background(), msg)

	if err != nil {
		t.Fatalf("expected error to be nil, but got %v", err)
	}
	if !store.SaveInvoked {
		t.Fatal("expected store.save() to be called")
	}

	if !acker.AckInvoked {
		t.Fatal("expected message.ack() to be called")
	}
}

func testShouldFailToUnmarshalBody(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *pkg.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return nil }
	body := []byte(`INVALID`)

	msg := message.New(acker, body)
	h := NewUserEmailChanged(store)
	err := h.Handle(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected to return err but got nil")
	}
	if store.SaveInvoked {
		t.Fatal("expected store.save() to not be called")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be called")
	}
}

func testHandleNotFoundError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *pkg.User) error { return pkg.ErrNotFound }
	acker.AckFunc = func(multiple bool) error {
		if multiple {
			t.Fatal("unexpected multiple")
		}
		return nil
	}
	body := []byte(`{
		"email": "foo@mail.com",
		"username": "foo",
		"status": "active"
	}`)

	msg := message.New(acker, body)
	h := NewUserEmailChanged(store)
	err := h.Handle(context.Background(), msg)
	if err != pkg.ErrNotFound {
		t.Fatalf("expected to return err but got %v", err)
	}
	if !store.SaveInvoked {
		t.Fatal("expected store.Save() to not be called")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be called")
	}
}

func testHandleUnexpectedSaveError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *pkg.User) error { return errors.New("unexpected error") }
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
		"status": "active"
	}`)

	msg := message.New(acker, body)
	h := NewUserEmailChanged(store)
	err := h.Handle(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected to return err but got nil")
	}
	if !store.SaveInvoked {
		t.Fatal("expected store.Save() to not be called")
	}
	if !acker.NackInvoked {
		t.Fatal("expected message.Ack() to be called")
	}
}

func testEmailChangeHandlerShouldFailToAck(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *pkg.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return errAcker }
	body := []byte(`INVALID`)

	msg := message.New(acker, body)
	h := NewUserEmailChanged(store)
	err := h.Handle(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected to return err but got nil")
	}
	if err != errAcker {
		t.Fatalf("expected to return errAcker, but got %v", err)
	}
	if store.SaveInvoked {
		t.Fatal("expected store.save() to not be called")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be called")
	}
}
