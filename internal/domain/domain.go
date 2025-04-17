package domain

import (
	"encoding/json"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type EventID string

func NewEventID() EventID { return EventID(NewID()) }
func NewID() string       { return gonanoid.Must(32) }

type EventBody any

type Event struct {
	ID          EventID    `json:"id"`
	PublishedAt *time.Time `json:"published_at"`
	Body        EventBody
}

type eventJSON struct {
	ID          EventID    `json:"id"`
	PublishedAt *time.Time `json:"published_at"`
	Body        json.RawMessage
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var tmp eventJSON
	err := json.Unmarshal(data, &tmp)
	e.ID = tmp.ID
	e.PublishedAt = tmp.PublishedAt
	if err == nil {
		e.Body, err = unmarshaller.UnmarshalEvent(tmp.Body)
	}
	return err
}

func NewDomainEvent(data EventBody) Event {
	return Event{ID: NewEventID(), Body: data}
}

type Unmarshaller interface {
	UnmarshalEvent([]byte) (EventBody, error)
}

var unmarshaller Unmarshaller

type UnmarshallerFunc func([]byte) (EventBody, error)

func (f UnmarshallerFunc) UnmarshalEvent(data []byte) (EventBody, error) { return f(data) }

func RegisterUnmarshaller(u Unmarshaller) {
	unmarshaller = u
}
