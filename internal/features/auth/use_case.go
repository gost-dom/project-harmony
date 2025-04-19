package auth

import "harmony/internal/domain"

// TODO: This is a general concept for domain logic - move to a general place

type Entity[T any] interface{ ID() T }

type DomainEvent = domain.Event

// UseCaseResult represents the outcome of a use case operating on a single
// entity or aggregate. The use case may result in an updated or new entity, as
// well as one or more domain events. Updating the entity and publishing the
// events must be an "atomic" operation. By atomic means, if the operation
// succeed, the use case must have been updated, and the events are guaranteed
// to be delivered in the future.
//
//   - If an entity has been updated, but events not published, an important
//     business operation may not execute.
//   - Publishing an event relating to an update that hasn't occurred may trigger
//     invalid business operations. E.g., they may read the current state which is
//     inconsistent with the event.
//
// For the second point, it is imperative that events are not published until
// AFTER a database transaction has committed.
type UseCaseResult[T any] struct {
	Entity T
	Events []DomainEvent
}

func UseCaseOfEntity[T any](e T) UseCaseResult[T] { return UseCaseResult[T]{Entity: e} }

func (useCase *UseCaseResult[T]) AddEvent(event DomainEvent) *UseCaseResult[T] {
	useCase.Events = append(useCase.Events, event)
	return useCase
}
