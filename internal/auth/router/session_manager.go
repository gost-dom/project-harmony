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

func (m *SessionManager) LoggedInUser(
	r *http.Request,
) (authAcc domain.AuthenticatedAccount, ok bool) {
	session, err := m.session(r)
	if err != nil {
		log.LogError(r.Context(), "SessionManager: load session error", err)
		return domain.AuthenticatedAccount{}, false
	}
	sessionValue, _ := session.Values[sessionAccountKey]
	accountID, ok := sessionValue.(domain.AccountID)
	if !ok {
		return domain.AuthenticatedAccount{}, false
	}
	acc, err := m.Repo.Get(r.Context(), accountID)
	if err != nil {
		log.LogError(r.Context(), "SessionManager: load account error", err)
		return domain.AuthenticatedAccount{}, false
	}
	authAcc, err = acc.Authenticated()
	ok = err == nil
	if err != nil {
		log.LogError(r.Context(), "SessionManager: Error createing authenticated account", err)
	}
	return
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
	deleteAllSessionValues(session)
	session.Values[sessionAccountKey] = account.ID
	// Prevent session fixation. Shouldn't be necessary, as we only store
	// one value in the session.
	return session.Save(req, w)
}

func (m SessionManager) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := m.session(r)
	if err != nil {
		return err
	}
	deleteAllSessionValues(session)
	expireCookie(session)
	return session.Save(r, w)
}

func (m SessionManager) session(r *http.Request) (*sessions.Session, error) {
	reg := sessions.GetRegistry(r)
	return reg.Get(m.SessionStore, sessionNameAuth)
}

// deleteAllSessionValues prevents leftover data from persisting when
//
// - Authenticating an already authenticated user
// - Logging out
//
// Although only an accountID is stored in the session, an aggressive strategy
// is chosen to mitigate security vulnerabilities in future versions of the
// code.
func deleteAllSessionValues(s *sessions.Session) {
	for k := range s.Values {
		delete(s.Values, k)
	}
}

func expireCookie(s *sessions.Session) {
	if s.Options != nil {
		s.Options.MaxAge = -1
	}
}
