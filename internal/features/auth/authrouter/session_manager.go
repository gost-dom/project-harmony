package authrouter

import (
	"harmony/internal/features/auth"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionNameAuth = "auth"
)

const sessionCookieName = "accountId"

type SessionManager struct {
	SessionStore sessions.Store
}

func (m *SessionManager) LoggedInUser(r *http.Request) *auth.Account {
	session, _ := m.SessionStore.Get(r, sessionNameAuth)
	if id, ok := session.Values[sessionCookieName]; ok {
		result := new(auth.Account)
		if strId, ok := id.(string); ok {
			result.Id = auth.AccountID(strId)
			return result
		}
	}
	return nil
}

func (s SessionManager) SetAccount(
	w http.ResponseWriter,
	req *http.Request,
	account auth.Account,
) error {
	session, err := s.SessionStore.Get(req, sessionNameAuth)
	if err != nil {
		return err
	}
	session.Values[sessionCookieName] = string(account.Id)
	return session.Save(req, w)
}
