package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/rafaeljesus/srv-consumer"
	"github.com/rafaeljesus/srv-consumer/mock"
	"github.com/rafaeljesus/srv-consumer/platform/message"
)

func TestUserStatusChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, *mock.UserStore, *mock.Acknowledger)
	}{
		{
			"when valid payload is supplied, then should successfully save user",
			testShouldSuccessfullyChangeUserStatus,
		},
		{
			"when invalid payload is supplied, then should fail to unmarshal body",
			testStatusChangeHandlerShouldFailToUnmarshalBody,
		},
		{
			"when Not found user is supplied, then should fail to save",
			testStatusChangeHandlerNotFoundError,
		},

		{
			"when unexpected error occurs, should be handled properly",
			testStatusChangeHandlerUnexpectedSaveError,
		},
		{
			"when unable to Ack message, error should be handled properly",
			testStatusChangeHandlerShouldFailToAck,
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

func testShouldSuccessfullyChangeUserStatus(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *srv.User) error {
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
	handler := NewUserStatusChanged(store)
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

func testStatusChangeHandlerShouldFailToUnmarshalBody(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *srv.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return nil }
	body := []byte(`INVALID`)

	msg := message.New(acker, body)
	h := NewUserStatusChanged(store)
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

func testStatusChangeHandlerNotFoundError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *srv.User) error { return srv.ErrNotFound }
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
	h := NewUserStatusChanged(store)
	err := h.Handle(context.Background(), msg)
	if err != srv.ErrNotFound {
		t.Fatalf("expected to return err but got %v", err)
	}
	if !store.SaveInvoked {
		t.Fatal("expected store.Save() to not be called")
	}
	if !acker.AckInvoked {
		t.Fatal("expected message.Ack() to be called")
	}
}

func testStatusChangeHandlerUnexpectedSaveError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *srv.User) error { return errors.New("unexpected error") }
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
	h := NewUserStatusChanged(store)
	err := h.Handle(context.Background(), msg)
	if err == nil {
		t.Fatalf("expected to return err but got nil")
	}
	if !store.SaveInvoked {
		t.Fatal("expected store.Save() to not be called")
	}
	if !acker.NackInvoked {
		t.Fatal("expected message.Nack() to be called")
	}
}

func testStatusChangeHandlerShouldFailToAck(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.SaveFunc = func(user *srv.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return errAcker }
	body := []byte(`INVALID`)

	msg := message.New(acker, body)
	h := NewUserStatusChanged(store)
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
