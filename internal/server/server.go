package server

import (
	"fmt"

	. "harmony/internal/features/auth/authrouter"
	"harmony/internal/project"
	"harmony/internal/server/views"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/a-h/templ"
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
	views.AuthLogin("/host", views.LoginFormData{}).Render(r.Context(), w)
}

func (s *Server) Init() {
	mux := http.NewServeMux()
	mux.Handle("/auth/", http.StripPrefix("/auth", s.AuthRouter))
	mux.Handle("GET /{$}", templ.Handler(views.Index()))
	mux.HandleFunc("GET /host/{$}", s.GetHost)
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
