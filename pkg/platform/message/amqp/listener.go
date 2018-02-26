package amqp

import (
	"context"
	"time"

	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
	"github.com/streadway/amqp"
)

type (
	// Stats expose methods for collecting metrics.
	Stats interface {
		// Start starts the timing metric.
		Start() time.Time
		// Track tracks operations for the given time.
		Track(t time.Time, err bool)
	}

	// Listener holds the fields for receiving incoming amqp messages.
	Listener struct {
		msgchan <-chan amqp.Delivery
		handler message.Handler
		stats   Stats
	}
)

// NewListener returns a configured listener.
func NewListener(key, ex string, c message.Consumer, h message.Handler, s Stats) (*Listener, error) {
	msgchan, err := c.Consume(key, ex)
	if err != nil {
		return nil, err
	}

	return &Listener{msgchan, h, s}, nil
}

// Run starts reading from amqp messages channel.
func (l *Listener) Run(ctx context.Context) error {
	for {
		select {
		case m, ok := <-l.msgchan:
			if ok {
				timing := l.stats.Start()
				msg := message.New(m, m.Body)
				err := l.handler.Handle(ctx, msg)
				l.stats.Track(timing, err == nil)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
