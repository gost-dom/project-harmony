package repotest

import (
	"context"
	"errors"
	"fmt"
	"harmony/internal/features/auth"
	"reflect"
	"testing"
)

var ErrDuplicateKey = errors.New("duplicate key")

type EntityTranslator[T, ID any] interface {
	ID(entity T) ID
}

type RepositoryStub[T any, ID comparable] struct {
	Translator EntityTranslator[T, ID]
	Entities   map[ID]*T
	Events     []auth.DomainEvent

	t testing.TB
}

func NewRepositoryStub[T any, ID comparable](
	t testing.TB,
	trans EntityTranslator[T, ID],
	entities ...*T,
) RepositoryStub[T, ID] {
	res := RepositoryStub[T, ID]{
		t:          t,
		Translator: trans,
		Entities:   make(map[ID]*T),
	}
	for _, e := range entities {
		res.Inject(e)
	}
	return res
}

// Inject is a test helper, allowing the test case to create an entity and
// inject a pointer, to the entity in the test case is updated; simplifying
// verification of state updates.
func (s *RepositoryStub[T, ID]) Inject(e *T) error {
	id := s.Translator.ID(*e)
	if _, exists := s.Entities[id]; exists {
		return ErrDuplicateKey
	}
	s.Entities[id] = e
	return nil
}

func (s *RepositoryStub[T, ID]) InsertEntity(_ context.Context, e T) error {
	ptr := new(T)
	*ptr = e
	return s.Inject(ptr)
}

func (s *RepositoryStub[T, ID]) Insert(ctx context.Context, e auth.UseCaseResult[T]) (T, error) {
	entity := e.Entity
	err := s.InsertEntity(ctx, entity)
	s.Events = append(s.Events, e.Events...)
	return entity, err
}

func (s RepositoryStub[T, ID]) Get(_ context.Context, id ID) (res T, err error) {
	if tmp, found := s.Entities[id]; found {
		res = *tmp
	} else {
		err = auth.ErrNotFound
	}
	return
}

func (s RepositoryStub[T, ID]) Update(_ context.Context, e T) (T, error) {
	id := s.Translator.ID(e)
	existing, ok := s.Entities[id]
	if !ok {
		var dummy T
		return dummy, auth.ErrNotFound
	}
	*existing = e
	return e, nil
}

func (s RepositoryStub[T, ID]) TestingT() testing.TB          { return s.t }
func (s RepositoryStub[T, ID]) AllEvents() []auth.DomainEvent { return s.Events }
func (s RepositoryStub[T, ID]) All() (res []*T) {
	res = make([]*T, len(s.Entities))
	i := 0
	for _, v := range s.Entities {
		res[i] = v
		i++
	}
	return res
}

func (repo RepositoryStub[T, ID]) Single() *T {
	ee := repo.All()
	if len(ee) != 1 {
		repo.t.Helper()
		repo.t.Errorf("repo.single: expected 1 element, had %d", len(ee))
		return nil
	}
	return ee[0]
}

func (repo RepositoryStub[T, ID]) Empty() bool { return len(repo.Entities) == 0 }

func (repo RepositoryStub[T, ID]) GetTestInstance(id ID) *T {
	if _, found := repo.Entities[id]; !found {
		panic(fmt.Sprintf("ID not found: %v", id))
	}
	return repo.Entities[id]
}

type E interface {
	TestingT() testing.TB
	AllEvents() []auth.DomainEvent
}

func SingleEventOfType[T any](e E) (res T) {
	t := e.TestingT()
	t.Helper()
	var found bool
	for _, ee := range e.AllEvents() {
		if r, ok := ee.Body.(T); ok {
			if found {
				t.Errorf("single-event: multiple instances of type %s", reflect.TypeFor[T]().Name())
			}
			res = r
			found = true
		}
	}
	if !found {
		t.Errorf("single-event: no event of type %s", reflect.TypeFor[T]().Name())
	}
	return
}
