package couchdb_test

import (
	"harmony/internal/couchdb"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
)

type Doc struct {
	Foo string
}

func TestDatabaseRoundtrip(t *testing.T) {
	conn := couchdb.DefaultConnection
	ctx := t.Context()

	// Insert a document
	id := gonanoid.Must()

	doc := Doc{Foo: "Bar"}
	rev, err := conn.Insert(ctx, id, doc)
	assert.NoError(t, err)
	assert.NotEmpty(t, rev, "A revision was returned")

	// Read the same doc
	var actual Doc
	rev, err = conn.Get(id, &actual)
	assert.NoError(t, err)

	// Verify they are equal
	assert.Equal(t, "Bar", actual.Foo)
	assert.Equal(t, doc, actual)

	actual.Foo = "Baz"
	_, err = conn.Update(ctx, id, rev, actual)
	assert.NoError(t, err, "Update error")

	var actualV2 Doc
	_, err = conn.Get(id, &actualV2)
	assert.NoError(t, err)
	assert.Equal(t, "Baz", actual.Foo)

	_, err = conn.Update(ctx, id, rev, actual)
	assert.ErrorIs(t, err, couchdb.ErrConflict)
}

func TestDatabaseBootstrap(t *testing.T) {
	if testing.Short() {
		// This isn't really a "slow" test, but it will try to connect to a
		// non-existing server - which could potentially have some timeout
		// issues in different environments.
		t.SkipNow()
	}
	_, err := couchdb.NewCouchConnection("http://invalid.localhost/")
	assert.ErrorIs(t, err, couchdb.ErrConn)
	assert.ErrorContains(t, err,
		"couchdb: connection error: ",
		"En error messages was appended to the standard error. Details not specified by the test")
}
