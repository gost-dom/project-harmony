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
	sessionCookieName = "accountId"
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
func (m *SessionManager) LoggedInUser(r *http.Request) (acc *domain.Account) {
	session, err := m.session(r)
	if err != nil {
		log.LogError(r.Context(), "SessionManager: load session error", err)
		return nil
	}
	if id, ok := session.Values[sessionCookieName]; ok {
		if id, ok := id.(domain.AccountID); ok {
			acc, err := m.Repo.Get(r.Context(), id)
			if err != nil {
				log.LogError(r.Context(), "SessionManager: load account error", err)
				return nil
			}
			return &acc
		}
	}
	return nil
}

func (m SessionManager) SetAccount(
	w http.ResponseWriter,
	req *http.Request,
	account domain.AuthenticatedAccount,
) error {
	reg := sessions.GetRegistry(req)
	session, err := reg.Get(m.SessionStore, sessionNameAuth)
	if err != nil {
		return err
	}
	session.Values[sessionCookieName] = account.ID
	return session.Save(req, w)
}

func (m SessionManager) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := m.session(r)
	if err != nil {
		return err
	}
	delete(session.Values, sessionCookieName)
	return session.Save(r, w)
}

func (m SessionManager) session(r *http.Request) (*sessions.Session, error) {
	return m.SessionStore.Get(r, sessionNameAuth)
}
