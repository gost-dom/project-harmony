package authrepo_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrCommunication = CommunicationError{}

type CommunicationError struct{ error }

func (CommunicationError) Is(err error) bool {
	_, is := err.(CommunicationError)
	return is
}
func (e CommunicationError) Error() string {
	return fmt.Sprintf("couchdb: communication error: %v", e.error)
}

type Doc struct {
	Foo string
}

type CouchConnection struct {
	url    string
	dbName string
}

// Bootstrap creates the database, as well as updates any design documents, such
// as views. Panics on unrecoverable errors, e.g., an invalid configuration.
func (c CouchConnection) Bootstrap() error {
	req, err := http.NewRequest("PUT", c.url, nil)
	if err != nil {
		// This is an unrecoverable error. The system configuration is wrong.
		// NewRequest will return an error is one of the arguments are invalid.
		// in this context that would mean that the URL is invalid.
		panic("couchdb: invalid configuration: url invalid")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CommunicationError{err}
	}
	switch resp.StatusCode {
	case 201, 202, 412:
		return nil
	case 400:
		// This is an unrecoverable error.
		panic("couchdb: invalid configuration: illegal database name")
	case 401:
		// This is an unrecoverable error.
		panic("couchdb: invalid configuration: bad credentials")
	default:
		panic("couchdb: unable to bootstrap database")
	}
}

func NewCouchConnection(couchURL string) (conn CouchConnection, err error) {
	var url *url.URL
	url, err = url.Parse(couchURL)
	conn = CouchConnection{couchURL, url.Path}
	if err == nil {
		err = conn.Bootstrap()
	}
	return
}

func (c CouchConnection) getDocId(location *url.URL) string {
	return strings.TrimPrefix(location.Path, c.dbName)
}

func (c CouchConnection) Insert(doc any) (id string, err error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	if err = enc.Encode(doc); err != nil {
		return
	}

	resp, err := http.Post(c.url, "application/json", &b)
	if err != nil {
		return
	}
	if resp.StatusCode != 201 {
		err = fmt.Errorf("couch: bad status code: %d", resp.StatusCode)
		return
	}
	loc, err := resp.Location()
	id = c.getDocId(loc)
	return
}

func (c CouchConnection) Get(id string, doc any) error {
	resp, err := http.Get(fmt.Sprintf("%s/%s", c.url, id))
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, &doc)
}

func TestDatabaseRoundtrip(t *testing.T) {
	conn, err := NewCouchConnection("http://admin:password@localhost:5984/harmony")
	assert.NoError(t, err)

	// Insert a document
	id, err := conn.Insert(Doc{Foo: "Bar"})
	assert.NoError(t, err)

	// Read the same doc
	var actual Doc
	err = conn.Get(id, &actual)
	assert.NoError(t, err)

	// Verify they are equal
	assert.Equal(t, "Bar", actual.Foo)
}

func TestDatabaseBootstrap(t *testing.T) {
	if testing.Short() {
		return
	}
	_, err := NewCouchConnection("http://invalid.localhost/")
	assert.ErrorIs(t, err, ErrCommunication)
}
