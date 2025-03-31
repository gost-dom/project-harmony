package repotest_test

import (
	"context"
	"harmony/internal/testing/repotest"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
)

type TestType struct{ ID string }

func NewTestType() TestType { return TestType{ID: gonanoid.Must()} }

type TestTypeTranslator struct{}

func (t TestTypeTranslator) ID(e TestType) string { return e.ID }

func TestFakeRepositoryDisallowsDuplicates(t *testing.T) {
	repo := repotest.NewRepositoryStub(t, TestTypeTranslator{})
	entity := TestType{ID: "DUMMY_ID"}
	err1 := repo.InsertEntity(context.Background(), entity)
	err2 := repo.InsertEntity(context.Background(), entity)
	assert.NoError(t, err1, "First insert")
	assert.ErrorIs(t, err2, repotest.ErrDuplicateKey)
}
