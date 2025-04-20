package messaging

import (
	"context"
	"harmony/internal/couchdb"
	"harmony/internal/domain"
	"harmony/internal/features/auth"
	"log/slog"
	"time"
)

type DomainEventUpdater interface {
	Update(context.Context, domain.Event) error
}

type MessageHandler struct {
	EventUpdater DomainEventUpdater
	Validator    auth.EmailValidator
}

func (h MessageHandler) ProcessDomainEvent(ctx context.Context, event domain.Event) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = h.Validator.ProcessDomainEvent(ctx, event); err == nil {
		event.MarkPublished()
		err = h.EventUpdater.Update(ctx, event)
	}
	return err
}

type MessagePump struct {
	couchdb.Connection
	Handler MessageHandler
}

func (h MessagePump) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting message pump")
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
