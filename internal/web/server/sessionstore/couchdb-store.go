package sessionstore

import (
	"context"
	"errors"
	"fmt"
	"harmony/internal/core"
	"harmony/internal/core/corerepo"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type CouchDBStore struct {
	db       *corerepo.Connection
	keyPairs [][]byte
	codecs   []securecookie.Codec
}

func NewCouchDBStore(db *corerepo.Connection, keyPairs [][]byte) CouchDBStore {
	return CouchDBStore{
		db,
		keyPairs,
		securecookie.CodecsFromPairs(keyPairs...),
	}
}

var _ sessions.Store = CouchDBStore{}

func (store CouchDBStore) Get(r *http.Request, name string) (s *sessions.Session, err error) {
	return sessions.GetRegistry(r).Get(store, name)
}

func (store CouchDBStore) docID(id string) string {
	return fmt.Sprintf("auth:sessions:%s", id)
}

func (store CouchDBStore) decodeSessionIdCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			err = nil
		}
		return "", err
	}
	return store.decodeIDCookie(name, c.Value)
}

func (store CouchDBStore) New(r *http.Request, name string) (session *sessions.Session, err error) {
	session = sessions.NewSession(store, name)
	*session.Options = store.opts()

	id, err := store.decodeSessionIdCookie(r, name)
	if err != nil || id == "" {
		return
	}

	var doc SessionDoc
	rev, err := store.db.Get(r.Context(), store.docID(id), &doc)
	if err != nil {
		if errors.Is(err, corerepo.ErrNotFound) {
			session.IsNew = true
			return session, nil
		}
		return
	}
	if err = store.decodeValues(doc.Values, session); err != nil {
		return
	}
	session.Values["_rev"] = rev
	session.ID = id
	session.IsNew = false
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

	doc := SessionDoc{
		ID:     session.ID,
		Values: v,
	}

	if session.ID == "" {
		session.ID = core.NewID()
		err = store.insert(r.Context(), session, doc)
	} else {
		if session.ID == "" {
			return fmt.Errorf("CouchDBStore.Save: session has no ID")
		}
		err = store.update(r.Context(), session, doc)
	}
	if err != nil {
		return err
	}

	encoded, err := store.encodeIDCookie(session)
	if err != nil {
		return fmt.Errorf("CouchDBStore.Save: encode cookie: %w", err)
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

func (store CouchDBStore) update(
	ctx context.Context, s *sessions.Session, doc SessionDoc,
) (err error) {
	rev, ok := s.Values["_rev"].(string)
	if !ok || rev == "" {
		return fmt.Errorf("CouchDBStore.update: session has no _rev value")
	}
	delete(s.Values, "_rev")
	s.Values["_rev"], err = store.db.Update(ctx, store.docID(s.ID), rev, doc)
	if err != nil {
		return fmt.Errorf("CouchDBStore.update: %w", err)
	}
	return
}

func (store CouchDBStore) insert(ctx context.Context, s *sessions.Session, doc SessionDoc) error {
	rev, err := store.db.Insert(ctx, store.docID(s.ID), doc)
	if err != nil {
		return fmt.Errorf("CouchDBStore.insert: %w", err)
	}
	s.Values["_rev"] = rev
	return nil
}

func (store CouchDBStore) decodeIDCookie(name, c string) (res string, err error) {
	err = securecookie.DecodeMulti(name, c, &res, store.codecs...)
	return
}

func (store CouchDBStore) encodeIDCookie(s *sessions.Session) (string, error) {
	return securecookie.EncodeMulti(s.Name(), s.ID, store.codecs...)
}

func (store CouchDBStore) encodeValues(s *sessions.Session) (string, error) {
	return securecookie.EncodeMulti(s.Name(), s.Values, store.codecs...)
}

func (store CouchDBStore) decodeValues(v string, s *sessions.Session) error {
	return securecookie.DecodeMulti(s.Name(), v, &s.Values, store.codecs...)
}

type SessionDoc struct {
	ID     string
	Values string
}
