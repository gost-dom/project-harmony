package sessionstore

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"
	"harmony/internal/couchdb"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type CouchDBStore struct {
	DB       *couchdb.Connection
	KeyPairs [][]byte
}

var _ sessions.Store = CouchDBStore{}

func (store CouchDBStore) Get(r *http.Request, name string) (s *sessions.Session, err error) {
	return sessions.GetRegistry(r).Get(store, name)
}

func (store CouchDBStore) docID(id string) string {
	return fmt.Sprintf("auth:sessions:%s", id)
}

func (store CouchDBStore) New(r *http.Request, name string) (s *sessions.Session, e error) {
	session := sessions.NewSession(store, name)
	*session.Options = store.opts()
	var err error
	cook, errCookie := r.Cookie(name)
	if errCookie == nil {
		id := decodeIDCookie(cook.Value)
		if id != "" {
			var doc SessionDoc
			rev, err := store.DB.Get(r.Context(), store.docID(id), &doc)
			if errors.Is(err, couchdb.ErrNotFound) {
				session.IsNew = true
				return session, nil
			}
			if err != nil {
				return nil, err
			}
			if err := store.decodeValues(doc.Values, session); err != nil {
				return nil, err
			}
			session.Values["_rev"] = rev
			session.ID = id
			session.IsNew = false
		}
	}
	return session, err
}

func (store CouchDBStore) opts() sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
	}
}

func (store CouchDBStore) Save(
	r *http.Request,
	w http.ResponseWriter,
	session *sessions.Session,
) (e error) {
	v, err := store.encodeValues(session)
	if err != nil {
		return fmt.Errorf("CouchDBStore.Save: encodeValues: %w", err)
	}
	// Set delete if max-age is < 0
	if session.Options.MaxAge <= 0 {
		// TODO: Delete in DB
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Generate a random session ID key suitable for storage in the DB
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32),
			), "=")
	}

	if session.ID == "" {
		return fmt.Errorf("CouchDBStore.Save: session has no ID")
	}
	rev, ok := session.Values["_rev"].(string)
	if !ok {
		err = store.insert(r.Context(), session, v)
	} else {
		delete(session.Values, "_rev")
		doc := SessionDoc{
			ID:     session.ID,
			Values: v,
		}
		session.Values["_rev"], err = store.DB.Update(r.Context(), store.docID(session.ID), rev, doc)
	}
	if err != nil {
		return err
	}

	encoded := encodeIDCookie(session)
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

func (store CouchDBStore) insert(ctx context.Context, s *sessions.Session, v string) error {
	doc := SessionDoc{
		ID:     s.ID,
		Values: v,
	}
	rev, err := store.DB.Insert(ctx, store.docID(s.ID), doc)
	if err != nil {
		return fmt.Errorf("CouchDBStore.insert: %w", err)
	}
	s.Values["_rev"] = rev
	return nil
}

// TODO: Encrypt the cookie.
func decodeIDCookie(c string) string            { return c }
func encodeIDCookie(s *sessions.Session) string { return s.ID }

func (store CouchDBStore) encodeValues(s *sessions.Session) (string, error) {
	codecs := securecookie.CodecsFromPairs(store.KeyPairs...)
	return securecookie.EncodeMulti(s.Name(), s.Values, codecs...)
}

func (store CouchDBStore) decodeValues(v string, s *sessions.Session) error {
	codecs := securecookie.CodecsFromPairs(store.KeyPairs...)
	return securecookie.DecodeMulti(s.Name(), v, &s.Values, codecs...)
}

type SessionDoc struct {
	ID     string
	Values string
}
