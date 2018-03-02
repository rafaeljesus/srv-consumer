package mock

import (
	"context"
	"sync"

	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
)

type (
	Handler struct {
		sync.RWMutex
		HandleInvoked bool
		HandleFunc    func(ctx context.Context, msg *message.Message) error
	}
)

func (h *Handler) Handle(ctx context.Context, msg *message.Message) error {
	h.Lock()
	defer h.Unlock()

	h.HandleInvoked = true
	return h.HandleFunc(ctx, msg)
}
