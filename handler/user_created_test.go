package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rafaeljesus/srv-consumer"
	"github.com/rafaeljesus/srv-consumer/handler"
	"github.com/rafaeljesus/srv-consumer/mock"
	"github.com/rafaeljesus/srv-consumer/platform/message"
	"github.com/rafaeljesus/srv-consumer/test_helper"
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
	h := handler.NewUserCreated(store)
	err := h.Handle(context.Background(), msg)

	testhelper.Ok(t, err)
	testhelper.Assert(t, store.AddInvoked != false, "expected store.Add() to be invoked")
	testhelper.Assert(t, acker.AckInvoked != false, "expected message.Ack() to be invoked")
}

func testFailToUnmarshalBody(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *srv.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return nil }
	body := []byte(``)

	msg := message.New(acker, body)
	h := handler.NewUserCreated(store)
	err := h.Handle(context.Background(), msg)

	testhelper.Assert(t, err != nil, "expected to return err: %v", err)
	testhelper.Assert(t, store.AddInvoked == false, "expected store.Add() to not be invoked")
	testhelper.Assert(t, acker.AckInvoked == true, "expected message.Ack() to be invoked")
}

func testFailAckWhenUnmarshalBodyError(t *testing.T, store *mock.UserStore, acker *mock.Acknowledger) {
	store.AddFunc = func(user *srv.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return errAcker }
	body := []byte(``)

	msg := message.New(acker, body)
	h := handler.NewUserCreated(store)
	err := h.Handle(context.Background(), msg)

	testhelper.Assert(t, err != nil, "expected to return err: %v", err)
	testhelper.Assert(t, store.AddInvoked == false, "expected store.Add() to not be invoked")
	testhelper.Assert(t, acker.AckInvoked == true, "expected message.Ack() to be invoked")
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
	h := handler.NewUserCreated(store)
	err := h.Handle(context.Background(), msg)

	testhelper.Assert(t, err == srv.ErrConflict, "expected to return ErrConflict but got %v", err)
	testhelper.Assert(t, store.AddInvoked == true, "expected store.Add() to not be invoked")
	testhelper.Assert(t, acker.AckInvoked == true, "expected message.Ack() to be invoked")
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
	h := handler.NewUserCreated(store)
	err := h.Handle(context.Background(), msg)

	testhelper.Assert(t, err != nil, "expected to return error ,but got : %v", err)
	testhelper.Assert(t, store.AddInvoked == true, "expected store.Add() to not be invoked")
	testhelper.Assert(t, acker.NackInvoked == true, "expected message.Nack() to be invoked")
}
