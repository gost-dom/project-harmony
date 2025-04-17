package domain

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func (e Event) MarshalJSON() ([]byte, error) {
	var js eventJSON
	typeName := types[reflect.TypeOf(e.Body)]

	if typeName == "" {
		return nil, fmt.Errorf("domain: Event.MarshalJSON: no registration for type %T", e.Body)
	}

	rawMessage, err := json.Marshal(e.Body)
	if err != nil {
		return nil, err
	}
	js.ID = e.ID
	js.PublishedAt = e.PublishedAt
	js.Type = typeName
	js.Body = (rawMessage)
	return json.Marshal(js)
}

type eventJSON struct {
	ID          EventID    `json:"id"`
	PublishedAt *time.Time `json:"published_at"`
	Type        string     `json:"type"`
	Body        json.RawMessage
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var tmp eventJSON
	err := json.Unmarshal(data, &tmp)
	e.ID = tmp.ID
	e.PublishedAt = tmp.PublishedAt
	if err == nil {
		val := reflect.New(names[tmp.Type])
		err = json.Unmarshal(tmp.Body, val.Interface())
		e.Body = EventBody(val.Elem().Interface())
	}
	return err
}

func NewDomainEvent(data EventBody) Event {
	return Event{ID: NewEventID(), Body: data}
}

type UnmarshallerFunc func([]byte) (EventBody, error)

func (f UnmarshallerFunc) UnmarshalEvent(data []byte) (EventBody, error) { return f(data) }

var types = make(map[reflect.Type]string)
var names = make(map[string]reflect.Type)

// Registers a unique name for an event type. The idiomatic name is the context name
// and type name explicing the "Event" suffix, separated by a dot. E.g.,
// AccountRegisteredEvent in the authentication area is
// "auth.AccountRegistered".
func RegisterType(typ_ reflect.Type, name string) {
	_, typeExists := types[typ_]
	_, nameExists := names[name]
	if typeExists {
		panic(fmt.Sprintf("domain: RegisterType: type already exists: %v", typ_))
	}
	if nameExists {
		panic(fmt.Sprintf("domain: RegisterType: name already exists: %s", name))
	}
	if typ_.Kind() != reflect.Struct {
		panic(fmt.Sprintf("domain: RegisterType: type must be a struct: %v", typ_))
	}

	types[typ_] = name
	names[name] = typ_
}
