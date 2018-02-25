package listener

import "github.com/streadway/amqp"

const (
	kind  = "topic"
	queue = "srv-consumer"
)

type (
	// Listener interpret (implement) Amqp consumer.
	Listener struct {
		ch *amqp.Channel
	}
)

// New returns a new Listener configured.
func New(ch *amqp.Channel) *Listener {
	return &Listener{ch}
}

// Listen creates a amqp consumer
func (l *Listener) Listen(exchange, key string) (<-chan amqp.Delivery, error) {
	if err := l.ch.ExchangeDeclare(exchange, kind, true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := l.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := l.ch.QueueBind(q.Name, key, exchange, false, nil); err != nil {
		return nil, err
	}

	return l.ch.Consume(q.Name, "", false, false, false, false, nil)
}
