package corerepo

import (
	"context"
	"harmony/internal/couchdb"
	"harmony/internal/domain"
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
