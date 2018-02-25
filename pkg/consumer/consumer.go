package consumer

import (
	"context"
	"time"

	"github.com/rafaeljesus/srv-consumer/pkg/message"
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

	// Listener is a method for binding routingKey, exchange and creating the amqp queue.
	Listener interface {
		Listen(routingKey, exchange string) (<-chan amqp.Delivery, error)
	}

	// Handler is the message handler.
	Handler interface {
		Handle(msg *message.Message) error
	}

	// Consumer holds the fields for receiving incoming amqp messages.
	Consumer struct {
		msgchan <-chan amqp.Delivery
		handler Handler
		stats   Stats
	}
)

// NewConsumer returns a configured consumer.
func New(key, ex string, l Listener, h Handler, s Stats) (*Consumer, error) {
	msgchan, err := l.Listen(key, ex)
	if err != nil {
		return nil, err
	}

	return &Consumer{msgchan, h, s}, nil
}

// Run starts reading from amqp messages channel.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		select {
		case m, ok := <-c.msgchan:
			if ok {
				timing := c.stats.Start()
				msg := message.New(m, m.Body)
				err := c.handler.Handle(msg)
				c.stats.Track(timing, err == nil)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
