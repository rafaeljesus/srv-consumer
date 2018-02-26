package message

type (
	// Acknowledger expose methods for acknowledge messages
	Acknowledger interface {
		// Ack acknowledge the message
		Ack(multiple bool) error
		// Nack negatively acknowledge the message
		Nack(multiple, requeue bool) error
		// Reject negatively acknowledge the message dropping the message if requeue if false
		Reject(requeue bool) error
	}

	// Message is the RabbitMQ message
	Message struct {
		Acknowledger
		Headers map[string]interface{}
		Body    []byte
	}
)

// New create new application message
func New(ac Acknowledger, body []byte) *Message {
	return &Message{
		Acknowledger: ac,
		Body:         body,
		Headers:      make(map[string]interface{}),
	}
}
