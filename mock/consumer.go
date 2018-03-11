package mock

import "github.com/streadway/amqp"

type (
	Consumer struct {
		ConsumeInvoked bool
		ConsumeFunc    func(routingKey, exchange string) (<-chan amqp.Delivery, error)
	}
)

func (c *Consumer) Consume(routingKey, exchange string) (<-chan amqp.Delivery, error) {
	c.ConsumeInvoked = true
	return c.ConsumeFunc(routingKey, exchange)
}
