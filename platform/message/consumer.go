package message

import (
	"context"

	"github.com/streadway/amqp"
)

type (
	// Consumer is a method for binding routingKey, exchange and creating the amqp queue.
	Consumer interface {
		Consume(routingKey, exchange string) (<-chan amqp.Delivery, error)
	}

	// Handler is the message handler.
	Handler interface {
		Handle(context context.Context, msg *Message) error
	}
)
