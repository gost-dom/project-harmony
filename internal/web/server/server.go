package server

import (
	"net/http"
	"path/filepath"

	authrouter "harmony/internal/auth/authrouter"
	hostrouter "harmony/internal/host/hostrouter"
	"harmony/internal/web"
	"harmony/internal/web/server/views"

	"github.com/a-h/templ"
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
		web.CSRFMiddleware,
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
