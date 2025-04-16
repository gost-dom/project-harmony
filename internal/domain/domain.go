package domain

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type EventID string

func NewEventID() EventID { return EventID(NewID()) }
func NewID() string       { return gonanoid.Must(32) }

type DomainEventData any

type DomainEvent struct {
	ID          EventID    `json:"id"`
	PublishedAt *time.Time `json:"published_at"`
	DomainEventData
}

func NewDomainEvent(data DomainEventData) DomainEvent {
	return DomainEvent{ID: NewEventID(), DomainEventData: data}
}
