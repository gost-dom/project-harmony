package router

import (
	"context"
	"harmony/internal/auth/domain"
	"harmony/internal/infrastructure/log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionNameAuth   = "auth"
	sessionAccountKey = "accountId"
)

type AccountGetter interface {
	Get(context.Context, domain.AccountID) (domain.Account, error)
}

type SessionManager struct {
	Repo         AccountGetter
	SessionStore sessions.Store
}

// TODO: Check account id is valid
// TODO: Return an AuthenticatedAccount if successful
func (m *SessionManager) LoggedInUser(r *http.Request) (*domain.Account, bool) {
	session, err := m.session(r)
	if err != nil {
		log.LogError(r.Context(), "SessionManager: load session error", err)
		return nil, false
	}
	sessionValue, _ := session.Values[sessionAccountKey]
	accountID, ok := sessionValue.(domain.AccountID)
	if !ok {
		return nil, false
	}
	acc, err := m.Repo.Get(r.Context(), accountID)
	if err != nil {
		log.LogError(r.Context(), "SessionManager: load account error", err)
		return nil, false
	}
	return &acc, true
}

func (m SessionManager) SetAccount(
	w http.ResponseWriter,
	req *http.Request,
	account domain.AuthenticatedAccount,
) error {
	session, err := m.session(req)
	if err != nil {
		return err
	}
	for k := range session.Values {
		// Prevent session fixation. Shouldn't be necessary, as we only store
		// one value in the session.
		delete(session.Values, k)
	}
	session.Values[sessionAccountKey] = account.ID
	return session.Save(req, w)
}

func (m SessionManager) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := m.session(r)
	if err != nil {
		return err
	}
	delete(session.Values, sessionAccountKey)
	return session.Save(r, w)
}

func (m SessionManager) session(r *http.Request) (*sessions.Session, error) {
	reg := sessions.GetRegistry(r)
	return reg.Get(m.SessionStore, sessionNameAuth)
}
