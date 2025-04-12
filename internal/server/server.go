package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"

	. "harmony/internal/features/auth/authrouter"
	"harmony/internal/project"
	"harmony/internal/server/views"

	"github.com/a-h/templ"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func staticFilesPath() string { return filepath.Join(project.Root(), "static") }

type Server struct {
	http.Handler
	SessionManager SessionManager
	AuthRouter     *AuthRouter
}

type sessionName string

func (s *Server) GetHost(w http.ResponseWriter, r *http.Request) {
	if account := s.SessionManager.LoggedInUser(r); account != nil {
		views.HostsPage().Render(r.Context(), w)
		return
	}
	// Not authenticated; show login page
	fmtNewLocation := fmt.Sprintf("/auth/login?redirectUrl=%s", url.QueryEscape("/hosts"))
	w.Header().Add("hx-push-url", fmtNewLocation)
	s.AuthRouter.RenderLogin(w, r)
}

func getCSRFCookie(r *http.Request) string {
	cookie, err := r.Cookie("csrftoken")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func CSRFProtection(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT", "PATCH", "DELETE":
			{
				r.ParseForm()
				formToken := r.FormValue("csrf-token")
				expectedToken := getCSRFCookie(r)

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

		token, err := gonanoid.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "csrftoken",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
			},
		)
		newCtx := context.WithValue(r.Context(), "csrf-token", token)
		newReq := r.WithContext(newCtx)
		h.ServeHTTP(w, newReq)
	})
}

func (s *Server) Init() {
	mux := http.NewServeMux()
	mux.Handle("/auth/", CSRFProtection(http.StripPrefix("/auth", s.AuthRouter)))
	mux.Handle("GET /{$}", CSRFProtection(templ.Handler(views.Index())))
	mux.Handle("GET /host/{$}", CSRFProtection(http.HandlerFunc(s.GetHost)))
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	s.Handler = noCache(mux)
}

func NewAuthRouter() *AuthRouter {
	res := &AuthRouter{}
	res.Init()
	return res
}

func New() *Server {
	res := &Server{
		AuthRouter: NewAuthRouter(),
	}
	res.Init()
	return res
}
