package handler

import (
	"context"
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
	store.SaveFunc = func(user *pkg.User) error { return nil }
	acker.AckFunc = func(multiple bool) error { return nil }

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
