package router

import (
	"harmony/internal/auth/domain"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionNameAuth   = "auth"
	sessionCookieName = "accountId"
)

type SessionManager struct {
	SessionStore sessions.Store
}

func (m SessionManager) session(r *http.Request) (*sessions.Session, error) {
	return m.SessionStore.Get(r, sessionNameAuth)
}

// TODO: Check account id is valid
func (m *SessionManager) LoggedInUser(r *http.Request) (acc *domain.Account) {
	session, _ := m.session(r)
	if id, ok := session.Values[sessionCookieName]; ok {
		result := new(domain.Account)
		if strId, ok := id.(string); ok && strId != "" {
			result.ID = domain.AccountID(strId)
			return result
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
	session.Values[sessionCookieName] = string(account.ID)
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
