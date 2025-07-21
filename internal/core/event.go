package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type EventID string

func NewEventID() EventID { return EventID(NewID()) }

// NewID generates a URL-friendly string identifier based on a cryptographically
// strong random number generator.
func NewID() string { return gonanoid.Must(32) }

type EventBody any

type DomainEvent struct {
	ID          EventID    `json:"id"`
	Rev         string     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at"`
	Body        EventBody
}

func (e *DomainEvent) MarkPublished() {
	if e.PublishedAt != nil {
		return
	}
	now := time.Now().UTC()
	e.PublishedAt = &now
}

func (e DomainEvent) MarshalJSON() ([]byte, error) {
	var js eventJSON
	typeName := types[reflect.TypeOf(e.Body)]

	if e.Body == nil {
		return nil, fmt.Errorf("domain: Event.MarshalJSON: body is nil")
	}
	if typeName == "" {
		return nil, fmt.Errorf("domain: Event.MarshalJSON: no registration for type %T", e.Body)
	}

	rawMessage, err := json.Marshal(e.Body)
	if err != nil {
		return nil, err
	}
	js.ID = e.ID
	js.CreatedAt = e.CreatedAt
	js.PublishedAt = e.PublishedAt
	js.Type = typeName
	js.Body = (rawMessage)
	return json.Marshal(js)
}

type eventJSON struct {
	ID          EventID    `json:"id"`
	Rev         string     `json:"_rev,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at"`
	Type        string     `json:"type"`
	Body        json.RawMessage
}

func (e *DomainEvent) UnmarshalJSON(data []byte) error {
	var rawEvent eventJSON
	if err := json.Unmarshal(data, &rawEvent); err != nil {
		return err
	}
	t := names[rawEvent.Type]
	if t == nil {
		return fmt.Errorf(
			"domain: Event.UnmarshalJSON: unknown message type: %s - json %s",
			rawEvent.Type, string(data),
		)
	}
	body := reflect.New(t)
	if err := json.Unmarshal(rawEvent.Body, body.Interface()); err != nil {
		return err
	}
	e.ID = rawEvent.ID
	e.Rev = rawEvent.Rev
	e.PublishedAt = rawEvent.PublishedAt
	e.CreatedAt = rawEvent.CreatedAt
	e.Body = EventBody(body.Elem().Interface())
	return nil
}

func NewDomainEvent(data EventBody) DomainEvent {
	return DomainEvent{
		ID: NewEventID(), Body: data,
		CreatedAt: time.Now().UTC(),
	}
}

type UnmarshallerFunc func([]byte) (EventBody, error)

func (f UnmarshallerFunc) UnmarshalEvent(data []byte) (EventBody, error) { return f(data) }

var types = make(map[reflect.Type]string)
var names = make(map[string]reflect.Type)

// Registers a unique name for an event type. The idiomatic name is the context name
// and type name explicing the "Event" suffix, separated by a dot. E.g.,
// AccountRegisteredEvent in the authentication area is
// "auth.AccountRegistered".
//
// Note, while the code could have inferred the name from the type itself, that
// would could data in the database to the code, preventing the ability to
// refactor.
func RegisterEventType(typ_ reflect.Type, name string) {
	_, typeExists := types[typ_]
	_, nameExists := names[name]
	if typ_.Kind() != reflect.Struct {
		panic(fmt.Sprintf("domain: RegisterType: type must be a struct: %v", typ_))
	}
	if typeExists {
		panic(fmt.Sprintf("domain: RegisterType: type already exists: %v", typ_))
	}
	if nameExists {
		panic(fmt.Sprintf("domain: RegisterType: name already exists: %s", name))
	}

	types[typ_] = name
	names[name] = typ_
}
