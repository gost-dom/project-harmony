package authrouter

import (
	"harmony/internal/features/auth/authdomain"
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

func (m *SessionManager) LoggedInUser(r *http.Request) *authdomain.Account {
	session, _ := m.SessionStore.Get(r, sessionNameAuth)
	if id, ok := session.Values[sessionCookieName]; ok {
		result := new(authdomain.Account)
		if strId, ok := id.(string); ok {
			result.ID = authdomain.AccountID(strId)
			return result
		}
	}
	return nil
}

func (s SessionManager) SetAccount(
	w http.ResponseWriter,
	req *http.Request,
	account authdomain.AuthenticatedAccount,
) error {
	session, err := s.SessionStore.Get(req, sessionNameAuth)
	if err != nil {
		return err
	}
	session.Values[sessionCookieName] = string(account.ID)
	return session.Save(req, w)
}
