package authrepo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/testing/domaintest"
	"io"
	"net/http"
	"net/url"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
)

// ErrConn indicates that an error occurred trying to communicate with CouchDB
// itself. Possible causes:
//   - Temporary condition such as a disconnected network
//   - Configuration issue, e.g., a wrong host name
var ErrConn = errors.New("couchdb: connection error")

var ErrConflict = errors.New("couchdb: conflict")
var ErrNotFound = errors.New("couchdb: not found")

func errUnexpectedStatusCode(resp *http.Response) error {
	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

type Doc struct {
	Foo string
}

type Document any

type couchDoc struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
	Document
}

type CouchConnection struct{ dbURL *url.URL }

// Bootstrap creates the database, as well as updates any design documents, such
// as views. Panics on unrecoverable errors, e.g., an invalid configuration.
func (c CouchConnection) Bootstrap() error {
	req, err := http.NewRequest("PUT", c.dbURL.String(), nil)
	if err != nil {
		// This is an unrecoverable error. The system configuration is wrong.
		// But this really shouldn't happen, because the only possible source of
		// error should be an invalid URL, but the url has already been parsed.
		panic("couchdb: invalid configuration: url invalid")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConn, err)
	}
	switch resp.StatusCode {
	case 201, 202, 412: // 412 means the database already exists.
		return nil
	case 400:
		// This is an unrecoverable error. The configured database name is not a
		// valid name in couchdb.
		panic("couchdb: invalid configuration: invalid database name")
	case 401:
		// This is an unrecoverable error. The configured credentials are wrong.
		panic("couchdb: invalid configuration: bad credentials")
	default:
		panic("couchdb: unable to bootstrap database")
	}
}

func NewCouchConnection(couchURL string) (conn CouchConnection, err error) {
	var url *url.URL
	url, err = url.Parse(couchURL)
	conn = CouchConnection{url}
	if err == nil {
		err = conn.Bootstrap()
	}
	return
}

// docUrl generates the full resource URL for couchDB.
func (c CouchConnection) docUrl(id string) string { return c.dbURL.JoinPath(id).String() }

// Insert creates a new document with the specified id. If the operation
// succeeds, the revision of the new document is returned in the rev return
// value.
func (c CouchConnection) Insert(id string, doc any) (rev string, err error) {
	if id == "" {
		err = errors.New("couchdb: missing id")
		return
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	if err = enc.Encode(doc); err != nil {
		return
	}

	url := c.docUrl(id)
	req, err := http.NewRequest("PUT", url, &b)
	if err != nil {
		err = fmt.Errorf("couchdb: put failed: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		rev = resp.Header.Get("Etag")
	case 409:
		err = ErrConflict
	default:
		err = fmt.Errorf("couch: insert id(%s): %w", id, errUnexpectedStatusCode(resp))
		return
	}
	return
}

func (c CouchConnection) Get(id string, doc any) (rev string, err error) {
	var resp *http.Response
	if resp, err = http.Get(c.docUrl(id)); err != nil {
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		cd := couchDoc{}
		if err = json.Unmarshal(bodyBytes, &cd); err == nil {
			err = json.Unmarshal(bodyBytes, &doc)
		}
		rev = cd.Rev
	case 404:
		err = ErrNotFound
	default:
		err = fmt.Errorf("couch: insert id(%s): %w", id, errUnexpectedStatusCode(resp))
	}
	return
}

// Update updates the document in the database. If successful, it will return
// the updated revision of the document. If there is a conflict, it will return
// ErrConflict
func (c CouchConnection) Update(id, oldRev string, doc any) (newRev string, err error) {
	var b bytes.Buffer
	var req *http.Request
	enc := json.NewEncoder(&b)
	if err = enc.Encode(doc); err != nil {
		return
	}
	req, err = http.NewRequest("PUT", c.docUrl(id), &b)
	if err != nil {
		err = fmt.Errorf("couchdb: put failed: %v", err)
		return
	}
	req.Header.Add("If-Match", oldRev)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrConn, err)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 201:
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		cd := couchDoc{}
		if err = json.Unmarshal(bodyBytes, &cd); err == nil {
			err = json.Unmarshal(bodyBytes, &doc)
		}
		newRev = cd.Rev
	case 409:
		err = ErrConflict
	default:
		err = fmt.Errorf("couchdb: unexpected http status code: %d", resp.StatusCode)
	}
	return
}

func TestDatabaseRoundtrip(t *testing.T) {
	conn, err := NewCouchConnection("http://admin:password@localhost:5984/harmony")
	assert.NoError(t, err)

	// Insert a document
	id := gonanoid.Must()

	doc := Doc{Foo: "Bar"}
	rev, err := conn.Insert(id, doc)
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
	_, err = conn.Update(id, rev, actual)
	assert.NoError(t, err, "Update error")

	var actualV2 Doc
	_, err = conn.Get(id, &actualV2)
	assert.NoError(t, err)
	assert.Equal(t, "Baz", actual.Foo)

	_, err = conn.Update(id, rev, actual)
	assert.ErrorIs(t, err, ErrConflict)
}

func TestDatabaseBootstrap(t *testing.T) {
	if testing.Short() {
		// This isn't really a "slow" test, but it will try to connect to a
		// non-existing server - which could potentially have some timeout
		// issues in different environments.
		t.SkipNow()
	}
	_, err := NewCouchConnection("http://invalid.localhost/")
	assert.ErrorIs(t, err, ErrConn)
	assert.ErrorContains(t, err,
		"couchdb: connection error: ",
		"En error messages was appended to the standard error. Details not specified by the test")
}

type accountEmailDoc struct {
	authdomain.AccountID
}

type AccountRepository struct {
	CouchConnection
}

func (r AccountRepository) accDocId(id authdomain.AccountID) string {
	return fmt.Sprintf("auth:account:%s", id)
}

func (r AccountRepository) addrDocId(addr string) string {
	return fmt.Sprintf("auth:account:email:%s", addr)
}
func (r AccountRepository) accEmailDocID(acc authdomain.Account) string {
	return r.addrDocId(acc.Email.Address)
}

func (r AccountRepository) Insert(acc authdomain.Account) error {
	_, err := r.CouchConnection.Insert(r.accDocId(acc.ID), acc)
	if err == nil {
		_, err = r.CouchConnection.Insert(
			r.accEmailDocID(acc),
			struct{ authdomain.AccountID }{acc.ID},
		)
	}
	return err
}

func (r AccountRepository) Get(id authdomain.AccountID) (res authdomain.Account, err error) {
	_, err = r.CouchConnection.Get(r.accDocId(id), &res)
	return
}

func (r AccountRepository) FindByEmail(email string) (res authdomain.Account, err error) {
	var doc accountEmailDoc
	_, err = r.CouchConnection.Get(r.addrDocId(email), &doc)
	if err == nil {
		res, err = r.Get(doc.AccountID)
	}
	return
}

func TestAccountRoundtrip(t *testing.T) {
	conn, err := NewCouchConnection("http://admin:password@localhost:5984/harmony")
	assert.NoError(t, err)
	repo := AccountRepository{conn}
	acc := domaintest.InitAccount()
	assert.NoError(t, repo.Insert(acc))
	reloaded, err := repo.Get(acc.ID)
	assert.Equal(t, acc, reloaded)

	foundByEmail, err := repo.FindByEmail(acc.Email.Address)
	assert.NoError(t, err, "Error finding by email")
	assert.Equal(t, acc, foundByEmail, "Entity found by email")
}

func TestDuplicateEmail(t *testing.T) {
	email := domaintest.NewAddress()
	acc1 := domaintest.InitAccount(domaintest.WithEmail(email))
	acc2 := domaintest.InitAccount(domaintest.WithEmail(email))
	conn, err := NewCouchConnection("http://admin:password@localhost:5984/harmony")
	assert.NoError(t, err)
	repo := AccountRepository{conn}
	assert.NoError(t, repo.Insert(acc1))
	assert.ErrorIs(t, repo.Insert(acc2), ErrConflict)
}
