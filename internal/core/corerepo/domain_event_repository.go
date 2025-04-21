package corerepo

import (
	"context"
	"encoding/json"
	"harmony/internal/couchdb"
	"harmony/internal/domain"
	"log/slog"
)

type DomainEventRepository struct {
	DB *couchdb.Connection
}

func (r DomainEventRepository) docID(e domain.Event) string {
	return "domain_event:" + string(e.ID)
}

func (r DomainEventRepository) Insert(ctx context.Context, e domain.Event) (domain.Event, error) {
	var err error
	e.Rev, err = r.DB.Insert(ctx, r.docID(e), e)
	return e, err
}

func (r DomainEventRepository) Update(ctx context.Context, e domain.Event) (domain.Event, error) {
	var err error
	e.Rev, err = r.DB.Update(ctx, r.docID(e), e.Rev, e)
	return e, err
}

// StreamOfEvents returns a channel of domain events. New events stored in the
// database will automatically be sent to the channel
func (r DomainEventRepository) StreamOfEvents(ctx context.Context) (<-chan domain.Event, error) {
	ch, err := r.DB.Changes(
		ctx,
		couchdb.ChangeOptFilter("events", "unpublished_domain_events"),
		couchdb.ChangeOptIncludeDocs(),
	)
	if err != nil {
		return nil, err
	}
	return domainEventsOfChangeEvents(ctx, ch), nil
}

// domainEventsOfChangeEvents takes a channel of CouchDB change events, assumed
// to contain new domain event documents, and transforms it to a channel of
// [domain.Event]
func domainEventsOfChangeEvents(
	ctx context.Context,
	ch <-chan couchdb.ChangeEvent,
) <-chan domain.Event {
	cha := make(chan domain.Event)
	go func() {
		defer close(cha)
		for changeEvent := range ch {
			var ev domain.Event
			err := json.Unmarshal(changeEvent.Doc, &ev)
			if err != nil {
				slog.ErrorContext(ctx, "corerepo: process event", "err", err)
				continue
			}
			cha <- ev
		}
	}()
	return cha
}
