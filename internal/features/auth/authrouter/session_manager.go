package authrouter

import (
	"fmt"
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
	reg := sessions.GetRegistry(r)
	session, _ := reg.Get(m.SessionStore, sessionNameAuth)
	fmt.Println("  **** GET USER", session.Values[sessionCookieName])
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
	reg := sessions.GetRegistry(req)
	session, err := reg.Get(s.SessionStore, sessionNameAuth)
	// session, err := s.SessionStore.Get(req, sessionNameAuth)
	if err != nil {
		return err
	}
	session.Values[sessionCookieName] = string(account.ID)
	fmt.Println("  **** SET ACCOUNT", session.Values[sessionCookieName])
	return session.Save(req, w)
}
