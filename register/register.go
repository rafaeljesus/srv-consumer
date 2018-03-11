package register

import (
	"context"
	"time"

	"github.com/rafaeljesus/srv-consumer/platform/message"
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

	// Register holds the fields for receiving incoming amqp messages.
	Register struct {
		msgchan <-chan amqp.Delivery
		handler message.Handler
		stats   Stats
	}
)

// New returns a configured register.
func New(key, ex string, c message.Consumer, h message.Handler, s Stats) (*Register, error) {
	msgchan, err := c.Consume(key, ex)
	if err != nil {
		return nil, err
	}

	return &Register{msgchan, h, s}, nil
}

// Run starts reading from amqp messages channel.
func (r *Register) Run(ctx context.Context) error {
	for {
		select {
		case m, ok := <-r.msgchan:
			if ok {
				timing := r.stats.Start()
				msg := message.New(m, m.Body)
				err := r.handler.Handle(ctx, msg)
				r.stats.Track(timing, err == nil)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
