package amqp

import "github.com/streadway/amqp"

const (
	kind  = "topic"
	queue = "srv-consumer"
)

type (
	// Consumer is a convenient way for binding exchange, queue and consumer.
	Consumer struct {
		ch *amqp.Channel
	}
)

// NewConsumer returns a new consumer configured.
func NewConsumer(ch *amqp.Channel) *Consumer {
	return &Consumer{ch}
}

// Consumer creates a amqp consumer
func (c *Consumer) Consume(key, exchange string) (<-chan amqp.Delivery, error) {
	if err := c.ch.ExchangeDeclare(exchange, kind, true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := c.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := c.ch.QueueBind(q.Name, key, exchange, false, nil); err != nil {
		return nil, err
	}

	return c.ch.Consume(q.Name, "", false, false, false, false, nil)
}
