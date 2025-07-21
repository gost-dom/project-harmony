package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	authrouter "harmony/internal/auth/authrouter"
	hostrouter "harmony/internal/host/hostrouter"
	"harmony/internal/web"
	"harmony/internal/web/server/views"

	"github.com/a-h/templ"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func staticFilesPath() string { return filepath.Join(projectRoot(), "static") }

type Server struct {
	http.Handler
	AuthMiddlewares authrouter.Middlewares
	AuthRouter      *authrouter.AuthRouter
	HostRouter      *hostrouter.HostRouter
}

type sessionName string

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

func CSRFProtection(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT", "PATCH", "DELETE":
			{
				r.ParseForm()
				formTokenId := r.FormValue("csrf-id")
				formToken := r.FormValue("csrf-token")
				expectedToken := getCSRFCookie(formTokenId, r)
				deleteCSRFCookie(formTokenId, w)
				if formToken != expectedToken || expectedToken == "" {
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

		var fn CSRFGenerator = func() (string, string) {
			id, err1 := gonanoid.New()
			token, err2 := gonanoid.New()
			if err := errors.Join(err1, err2); err != nil {
				slog.Error("Error generating token", "err", err)
				return "", ""
			}
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
			return id, token
		}

		web.SetReqValue(&r, web.CtxKeyCSRFTokenSrc, fn)
		h.ServeHTTP(w, r)
	})
}

type CSRFGenerator = func() (string, string)

// Init implements interface [surgeon.Initer].
func (s *Server) Init() {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(views.Index()))
	mux.Handle("/auth/", http.StripPrefix("/auth", s.AuthRouter))
	mux.Handle("GET /host", authrouter.RequireAuth(s.HostRouter.Index()))
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	s.Handler = web.ApplyMiddlewares(mux,
		web.Logger,
		web.MiddlewareFunc(noCache),
		web.MiddlewareFunc(CSRFProtection),
		s.AuthMiddlewares,
	)
}

func New() *Server {
	res := &Server{
		AuthRouter: authrouter.New(),
		HostRouter: hostrouter.New(),
	}
	res.Init()
	return res
}
