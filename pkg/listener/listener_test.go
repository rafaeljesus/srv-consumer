package listener

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rafaeljesus/srv-consumer/pkg/mock"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
	"github.com/streadway/amqp"
)

var (
	amqpError = errors.New("amqp error")
)

func TestListener(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *mock.Consumer, *mock.Handler, *mock.Stats)
	}{
		{
			"create new listener",
			testCreateNewListener,
		},
		{
			"fail to create new listener",
			testFailToCreateNewListener,
		},
		{
			"run listener",
			testRunListener,
		},
		{
			"handle context done",
			testHandlerContextDone,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			c := new(mock.Consumer)
			h := new(mock.Handler)
			s := new(mock.Stats)
			test.function(t, c, h, s)
		})
	}
}

func testCreateNewListener(t *testing.T, consumer *mock.Consumer, handler *mock.Handler, stats *mock.Stats) {
	consumer.ConsumeFunc = func(key, ex string) (<-chan amqp.Delivery, error) {
		if key == "" {
			t.Fatalf("unexpected routingKey: %s", key)
		}
		if ex == "" {
			t.Fatalf("unexpected exchange: %s", ex)
		}
		return make(<-chan amqp.Delivery), nil
	}

	if _, err := New("key", "ex", consumer, handler, stats); err != nil {
		t.Fatalf("expected to create new listener: %v", err)
	}
	if !consumer.ConsumeInvoked {
		t.Fatal("expected consumer.Consume() to be invoked")
	}
}

func testFailToCreateNewListener(t *testing.T, consumer *mock.Consumer, handler *mock.Handler, stats *mock.Stats) {
	consumer.ConsumeFunc = func(key, ex string) (<-chan amqp.Delivery, error) { return nil, amqpError }
	if _, err := New("key", "ex", consumer, handler, stats); err != amqpError {
		t.Fatalf("expected to have amqpError: %v", err)
	}
	if !consumer.ConsumeInvoked {
		t.Fatal("expected consumer.Consume() to be invoked")
	}
}

func testRunListener(t *testing.T, consumer *mock.Consumer, handler *mock.Handler, stats *mock.Stats) {
	msgchan := make(chan amqp.Delivery)
	consumer.ConsumeFunc = func(key, ex string) (<-chan amqp.Delivery, error) {
		if key == "" {
			t.Fatalf("unexpected routingKey: %s", key)
		}
		if ex == "" {
			t.Fatalf("unexpected exchange: %s", ex)
		}
		return msgchan, nil
	}
	handler.HandleFunc = func(ctx context.Context, m *message.Message) error {
		if m == nil {
			t.Fatalf("unexpected message: %v", m)
		}
		return nil
	}
	stats.TrackFunc = func(tm time.Time, ok bool) {
		if tm.IsZero() {
			t.Fatal("unexpected time value")
		}
		if !ok {
			t.Fatalf("unexpected ok: %t", ok)
		}
	}

	l, err := New("key", "ex", consumer, handler, stats)
	if err != nil {
		t.Fatalf("expected to create new listener: %v", err)
	}
	if !consumer.ConsumeInvoked {
		t.Fatal("expected consumer.Consume() to be invoked")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go l.Run(ctx)
	msgchan <- amqp.Delivery{Body: []byte(`foo`)}

	handler.RLock()
	defer handler.RUnlock()
	if !handler.HandleInvoked {
		t.Fatal("expected handle.Handler() to be invoked")
	}

	stats.RLock()
	defer stats.RUnlock()
	if !stats.TrackInvoked {
		t.Fatal("expected stats.Track() to be invoked")
	}
}

func testHandlerContextDone(t *testing.T, consumer *mock.Consumer, handler *mock.Handler, stats *mock.Stats) {
	consumer.ConsumeFunc = func(key, ex string) (<-chan amqp.Delivery, error) { return make(chan amqp.Delivery), nil }
	handler.HandleFunc = func(ctx context.Context, m *message.Message) error { return nil }
	stats.TrackFunc = func(tm time.Time, ok bool) {}

	l, err := New("key", "ex", consumer, handler, stats)
	if err != nil {
		t.Fatalf("expected to create new listener: %v", err)
	}
	if !consumer.ConsumeInvoked {
		t.Fatal("expected consumer.Consume() to be invoked")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go l.Run(ctx)
	cancel()

	if handler.HandleInvoked {
		t.Fatal("expected handle.Handler() to not be invoked")
	}
	if stats.TrackInvoked {
		t.Fatal("expected stats.Track() to not be invoked")
	}
}
