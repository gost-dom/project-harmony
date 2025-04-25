package server

import (
	"context"
	"errors"
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

type StatusRecorder struct {
	http.ResponseWriter
	Code int
}

func (r *StatusRecorder) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
	r.Code = code
}

func log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{ResponseWriter: w}
		fmt.Printf("HTTP Request. Method: %s - Path: %s\n", r.Method, r.URL.Path)
		fmt.Println("Rewritten", r.Context().Value("rewritten"))
		h.ServeHTTP(rec, r)
		fmt.Printf("HTTP Resp: %d\n", rec.Code)
	})
}

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
	fmtNewLocation := fmt.Sprintf("/auth/login?redirectUrl=%s", url.QueryEscape("/host"))
	w.Header().Add("hx-push-url", fmtNewLocation)
	s.AuthRouter.RenderHost(w, r)
}

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

func rewriter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rewriter http.Handler = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				r = r.WithContext(context.WithValue(r.Context(), "rewritten", true))
				fmt.Println("***  Rewrite", r.URL)
				h.ServeHTTP(w, r)
			})
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "rewriter", rewriter)))
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
		newCtx := context.WithValue(r.Context(), "tokenSource", fn)
		newReq := r.WithContext(newCtx)
		h.ServeHTTP(w, newReq)
	})
}

type CSRFGenerator = func() (string, string)

func (s *Server) Init() {
	mux := http.NewServeMux()
	mux.Handle("/auth/", http.StripPrefix("/auth", s.AuthRouter))
	mux.Handle("GET /{$}", templ.Handler(views.Index()))
	mux.Handle("GET /host", http.HandlerFunc(s.GetHost))
	mux.Handle("GET /host/", http.HandlerFunc(s.GetHost))
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	s.Handler = rewriter(log(noCache(CSRFProtection(
		mux))))
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
