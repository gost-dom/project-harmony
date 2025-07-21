package messaging

import (
	"context"
	"harmony/internal/core"
	"harmony/internal/core/corerepo"
	"harmony/internal/features/auth"
	"log/slog"
	"time"
)

type DomainEventUpdater interface {
	Update(context.Context, core.DomainEvent) (core.DomainEvent, error)
}

type MessageHandler struct {
	EventUpdater DomainEventUpdater
	Validator    *auth.EmailValidator
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		nil,
		auth.NewEmailValidator(),
	}
}

func (h MessageHandler) ProcessDomainEvent(ctx context.Context, event core.DomainEvent) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = h.Validator.ProcessDomainEvent(ctx, event); err == nil {
		event.MarkPublished()
		_, err = h.EventUpdater.Update(ctx, event)
	}
	return err
}

type MessagePump struct {
	corerepo.MessageSource
	corerepo.DomainEventRepository
	Handler MessageHandler
}

func (h MessagePump) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting message pump")
	if ctx == nil {
		ctx = context.Background()
	}
	err := h.MessageSource.StartListener(ctx)
	if err != nil {
		return err
	}
	ch, err := h.DomainEventRepository.StreamOfEvents(ctx)
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
