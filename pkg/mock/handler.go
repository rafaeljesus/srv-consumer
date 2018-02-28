package mock

import (
	"context"

	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
)

type (
	Handler struct {
		HandleInvoked bool
		HandleFunc    func(ctx context.Context, msg *message.Message) error
	}
)

func (h *Handler) Handle(ctx context.Context, msg *message.Message) error {
	h.HandleInvoked = true
	return h.HandleFunc(ctx, msg)
}
