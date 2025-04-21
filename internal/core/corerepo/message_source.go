package corerepo

import (
	"context"
	"encoding/json"
	"errors"
	"harmony/internal/couchdb"
	"harmony/internal/domain"
	"log/slog"
)

type MessageSource struct {
	DomainEventRepository
	DB *couchdb.Connection
}

type DocumentWithEvents[T any] struct {
	ID       string         `json:"_id,omitempty"`
	Rev      string         `json:"_rev,omitempty"`
	Document T              `json:"doc"`
	Events   []domain.Event `json:"events,omitempty"`
}

func (c MessageSource) StartListener(
	ctx context.Context,
) (err error) {
	slog.InfoContext(ctx, "corerepo: Connection.StartListener")
	return c.processNewDomainEvents(ctx)
}

// processNewDomainEvents collects domain events from entity documents, and
// extracts them to dedicated domain event documents
func (c MessageSource) processNewDomainEvents(ctx context.Context) (err error) {
	ch, err := c.DB.Changes(
		ctx,
		couchdb.ChangeOptFilter("events", "aggregate_events"),
		couchdb.ChangeOptIncludeDocs(),
	)
	if err != nil {
		return
	}
	go func() {
		for doc := range getNewEntityEvents(ctx, ch) {
			c.processNewEntity(ctx, doc)
		}
	}()
	return nil
}

func (c MessageSource) processNewEntity(
	ctx context.Context,
	doc DocumentWithEvents[json.RawMessage],
) {
	for _, domainEvent := range doc.Events {
		_, err := c.DomainEventRepository.Insert(ctx, domainEvent)
		if err != nil && !errors.Is(err, couchdb.ErrConflict) {
			slog.ErrorContext(ctx, "corerepo: insert domain event", "err", err)
			return
		}
	}
	doc.Events = nil
	_, err := c.DB.Update(ctx, doc.ID, doc.Rev, doc)
	if err != nil {
		slog.ErrorContext(ctx, "corerepo: process event", "err", err)
		return
	}
}

func getNewEntityEvents(
	ctx context.Context,
	ch <-chan couchdb.ChangeEvent,
) <-chan DocumentWithEvents[json.RawMessage] {
	res := make(chan DocumentWithEvents[json.RawMessage])
	go func() {
		defer close(res)
		for changeEvent := range ch {
			var doc DocumentWithEvents[json.RawMessage]
			err := json.Unmarshal(changeEvent.Doc, &doc)
			if err != nil {
				slog.ErrorContext(ctx, "corerepo: process event document", "err", err)
				continue
			}
			res <- doc
		}
	}()
	return res
}
