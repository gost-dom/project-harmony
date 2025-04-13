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

type Doc struct {
	Foo string
}

type CouchConnection struct {
	url    string
	dbName string
}

func NewCouchConnection(couchURL string) (CouchConnection, error) {
	url, err := url.Parse(couchURL)
	dbName := url.Path
	return CouchConnection{couchURL, dbName}, err
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
		err = fmt.Errorf("CouchConnection.Insert: Unexpected status code: %d", resp.StatusCode)
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
