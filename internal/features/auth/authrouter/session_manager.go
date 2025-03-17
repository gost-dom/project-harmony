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
			result.Id = auth.AccountId(strId)
			return result
		}
	}
	return nil
}
