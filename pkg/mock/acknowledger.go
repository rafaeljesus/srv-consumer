package mock

type (
	Acknowledger struct {
		AckInvoked bool
		AckFunc    func(multiple bool) error

		NackInvoked bool
		NackFunc    func(multiple, requeue bool) error

		RejectInvoked bool
		RejectFunc    func(requeue bool) error
	}
)

func (c *Acknowledger) Ack(multiple bool) error {
	c.AckInvoked = true
	return c.AckFunc(multiple)
}

func (c *Acknowledger) Nack(multiple, requeue bool) error {
	c.NackInvoked = true
	return c.NackFunc(multiple, requeue)
}

func (c *Acknowledger) Reject(requeue bool) error {
	c.RejectInvoked = true
	return c.RejectFunc(requeue)
}
