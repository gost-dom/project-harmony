package messaging

import (
	"context"
	"harmony/internal/couchdb"
	"harmony/internal/domain"
	"harmony/internal/features/auth"
	"log/slog"
	"time"
)

type MessageHandler struct {
	Validator auth.EmailValidator
}

func (h MessageHandler) ProcessDomainEvent(ctx context.Context, event domain.Event) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	return h.Validator.ProcessDomainEvent(ctx, event)
}

type MessagePump struct {
	couchdb.Connection
	Handler MessageHandler
}

func (h MessagePump) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ch, err := h.Connection.StartListener(ctx)
	if err != nil {
		return err
	}
	go func() {
		for event := range ch {
			if err := h.Handler.ProcessDomainEvent(ctx, event); err != nil {
				slog.ErrorContext(ctx, "MessageHandler: error processing", "err", err)
			}
		}
	}()
	return nil
}
