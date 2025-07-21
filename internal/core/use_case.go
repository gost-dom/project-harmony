package core

// UseCaseResult represents the outcome of a "use case" operating on a single
// entity or aggregate.
//
// A use case may result in an updated entity, as well as a collection of
// "domain events". Wrapping this in a specific type helps implement the
// "transactional output" pattern.
//
// This type makes it decoupled from database technology, e.g., a relational
// database will wrap multiple writes in a transaction, where a document
// database typically have a single document as the boundary of consistency, so
// events must be stored _with_ the entity; and later processed.
//
// Note: Some programs take the approach to add events _to_ the entity type.
// The applications where I've seen this, it has been due to a technical
// limitation of persistence libraries, e.g. ORMs. The domain events are
// conceptually _not_ part of the entity; the entity is the _source_ of events.
type UseCaseResult[T any] struct {
	Entity T
	Events []DomainEvent
}

func UseCaseOfEntity[T any](e T) UseCaseResult[T] { return UseCaseResult[T]{Entity: e} }

func (useCase *UseCaseResult[T]) AddEvent(event DomainEvent) *UseCaseResult[T] {
	useCase.Events = append(useCase.Events, event)
	return useCase
}
