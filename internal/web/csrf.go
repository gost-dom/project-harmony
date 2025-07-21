package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func csrfCookieName(id string) string { return fmt.Sprintf("csrf-%s", id) }

func getCSRFCookie(id string, r *http.Request) string {
	cookie, err := r.Cookie(csrfCookieName(id))
	if err != nil {
		return ""
	}
	return cookie.Value
}

func deleteCSRFCookie(id string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName(id),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})
}

func csrfMiddleware(h http.Handler) http.Handler {
	// The CSRF middleware will generate a new cookie value on each request.
	// If a page renders a form, it must include the same value. Upon receiving
	// a mutating request, a form must be submitted containing values matching
	// the cookie.
	//
	// This scheme would ruin it for users opening multiple windows, so each
	// request generates a unique cookie name; also deleting the previous value.
	//
	// Cookies have a max-age of 1 hour. It is sufficient for type types of form
	// that the system is currently serving.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT", "PATCH", "DELETE":
			// For a mutating request, verify that the request body contains a
			// form with values matching those in a cookie.
			{
				r.ParseForm()
				formTokenId := r.FormValue("csrf-id")
				formToken := r.FormValue("csrf-token")
				expectedToken := getCSRFCookie(formTokenId, r)
				deleteCSRFCookie(formTokenId, w)
				if expectedToken == "" {
					http.Error(w, "Missing CSRF token", http.StatusBadRequest)
					return
				}
				if formToken != expectedToken {
					http.Error(w, "Invalid CSRF token", http.StatusBadRequest)
					return
				}
			}
		case "GET", "TRACE", "OPTIONS", "HEAD":
			// These methods are not mutating, so CSRF protection is not necessary
		default:
			// Unexpected HTTP method. If it's unexpected, it's unlikely that it
			// will have an effect on the server, so a warning log message
			// should be fine.
			slog.WarnContext(
				r.Context(),
				"CSRFMiddleware: Unexpected HTTP Method",
				"method",
				r.Method,
			)
		}

		var fn csrfGenerator = func() CSRFFields {
			id := gonanoid.Must()
			token := gonanoid.Must()
			http.SetCookie(
				w,
				&http.Cookie{
					Name:     csrfCookieName(id),
					Value:    token,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
					MaxAge:   3600,
				},
			)
			return CSRFFields{id, token}
		}

		SetReqValue(&r, CtxKeyCSRFTokenSrc, fn)
		h.ServeHTTP(w, r)
	})
}

type csrfGenerator = func() CSRFFields

// CSRFFields contains a pair of ID and Token value.
type CSRFFields struct {
	ID    string
	Token string
}

// GetCSRFFields generates a set of CSRF id/token values, and saves them in a
// cookie.
func GetCSRFFields(ctx context.Context) (fields CSRFFields, ok bool) {
	var g csrfGenerator
	g, ok = ctx.Value(CtxKeyCSRFTokenSrc).(csrfGenerator)
	if ok {
		fields = g()
	}
	return
}

var CSRFMiddleware = MiddlewareFunc(csrfMiddleware)
