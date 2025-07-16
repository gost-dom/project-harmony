package ctx

import "context"

const (
	AuthLoggedIn = "auth:loggedIn"
)

func IsLoggedIn(c context.Context) bool {
	if val, ok := c.Value(AuthLoggedIn).(bool); ok {
		return val
	}
	return false
}

func SetIsLoggedIn(c context.Context, v bool) context.Context {
	return context.WithValue(c, AuthLoggedIn, v)
}
