package corerepo

import (
	"context"
	"encoding/json"
	"harmony/internal/core"
	"harmony/internal/infrastructure/log"
	"net/url"
)

type DomainEventRepository struct {
	DB *Connection
}

func (r DomainEventRepository) docID(e core.DomainEvent) string {
	return "domain_event:" + string(e.ID)
}

func (r DomainEventRepository) Insert(
	ctx context.Context,
	e core.DomainEvent,
) (core.DomainEvent, error) {
	var err error
	e.Rev, err = r.DB.Insert(ctx, r.docID(e), e)
	return e, err
}

func (r DomainEventRepository) Update(
	ctx context.Context,
	e core.DomainEvent,
) (core.DomainEvent, error) {
	var err error
	e.Rev, err = r.DB.Update(ctx, r.docID(e), e.Rev, e)
	return e, err
}

// StreamOfEvents returns a channel of domain events. New events stored in the
// database will automatically be sent to the channel
func (r DomainEventRepository) StreamOfEvents(
	ctx context.Context,
) (<-chan core.DomainEvent, error) {
	ch, err := r.DB.Changes(
		ctx,
		ChangeOptViewFilter("events", "unpublished_domain_events"),
		ChangeOptIncludeDocs(),
	)
	if err != nil {
		return nil, err
	}
	return r.domainEventsOfChangeEvents(ctx, ch)
}

func (r DomainEventRepository) getCurrentDomainEvents() ([]core.DomainEvent, error) {
	v := make(url.Values)
	v.Add("include_docs", "true")
	var res DocsViewResult[core.DomainEvent]
	_, err := r.DB.GetPath("_design/events/_view/unpublished_domain_events", v, &res)
	return res.Docs(), err
}

// domainEventsOfChangeEvents takes a channel of CouchDB change events, assumed
// to contain new domain event documents, and transforms it to a channel of
// [core.DomainEvent]
func (r DomainEventRepository) domainEventsOfChangeEvents(
	ctx context.Context,
	ch <-chan ChangeEvent,
) (<-chan core.DomainEvent, error) {
	events, err := r.getCurrentDomainEvents()
	if err != nil {
		return nil, err
	}

	cha := make(chan core.DomainEvent, DEFAULT_EVENT_BUFFER_SIZE)
	go func() {
		defer close(cha)

		for _, e := range events {
			select {
			case cha <- e:
			case <-ctx.Done():
				return
			}
		}

		for changeEvent := range ch {
			var ev core.DomainEvent
			err := json.Unmarshal(changeEvent.Doc, &ev)
			if err != nil {
				log.Error(ctx, "corerepo: process event", "err", err)
				continue
			}
			select {
			case cha <- ev:
			case <-ctx.Done():
				return
			}
		}
	}()
	return cha, nil
}
