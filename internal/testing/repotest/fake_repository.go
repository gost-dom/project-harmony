package repotest

import (
	"context"
	"errors"
	"harmony/internal/features/auth"
	"reflect"
	"testing"
)

var ErrDuplicateKey = errors.New("duplicate key")

type EntityTranslator[T any] interface {
	ID(entity T) string
}

type RepositoryStub[T any] struct {
	Translator EntityTranslator[T]
	Entities   map[string]*T
	Events     []auth.DomainEvent

	t testing.TB
}

func NewRepositoryStub[T any](t testing.TB, trans EntityTranslator[T]) RepositoryStub[T] {
	return RepositoryStub[T]{
		t:          t,
		Translator: trans,
		Entities:   make(map[string]*T),
	}
}

func (s *RepositoryStub[T]) InsertEntity(_ context.Context, e T) error {
	id := s.Translator.ID(e)
	if _, exists := s.Entities[id]; exists {
		return ErrDuplicateKey
	}
	s.Entities[id] = new(T)
	*s.Entities[id] = e
	return nil
}

func (s *RepositoryStub[T]) Insert(ctx context.Context, e auth.UseCaseResult[T]) error {
	entity := e.Entity
	err := s.InsertEntity(ctx, entity)
	s.Events = append(s.Events, e.Events...)
	return err
}

func (s RepositoryStub[T]) Get(_ context.Context, id string) (res T) {
	tmp, found := s.Entities[id]
	if found {
		res = *tmp
	} else {
		s.t.Helper()
		s.t.Errorf("repo.get: id not found: %v", id)
		tmp = new(T)
	}
	return
}

func (s RepositoryStub[T]) TestingT() testing.TB          { return s.t }
func (s RepositoryStub[T]) AllEvents() []auth.DomainEvent { return s.Events }
func (s RepositoryStub[T]) All() (res []*T) {
	res = make([]*T, len(s.Entities))
	i := 0
	for _, v := range s.Entities {
		res[i] = v
		i++
	}
	return res
}

func (repo RepositoryStub[T]) Single() *T {
	ee := repo.All()
	if len(ee) != 1 {
		repo.t.Helper()
		repo.t.Errorf("repo.single: expected 1 element, had %d", len(ee))
		return nil
	}
	return ee[0]
}

func (repo RepositoryStub[T]) Empty() bool { return len(repo.Entities) == 0 }

func (repo RepositoryStub[T]) GetTestInstance(id string) *T {
	if _, found := repo.Entities[id]; !found {
		panic("ID not found: " + id)
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
