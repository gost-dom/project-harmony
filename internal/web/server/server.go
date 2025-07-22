package server

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	authrouter "harmony/internal/auth/router"
	hostrouter "harmony/internal/host/router"
	"harmony/internal/web"
	"harmony/internal/web/server/views"

	"github.com/a-h/templ"
)

var noCache = web.MiddlewareFunc(func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
})

func staticFilesPath() string { return filepath.Join(projectRoot(), "static") }

type Server struct {
	http.Handler
	AuthMiddlewares authrouter.Middlewares
	AuthRouter      *authrouter.AuthRouter
	HostRouter      *hostrouter.HostRouter
}

func stripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if p == "" {
			p = "/"
			rp = "/"
		}
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}

func subRoute(mux *http.ServeMux, path string, mw web.MiddlewareFunc, h http.Handler) {
	h = mw(stripPrefix(path, h))
	mux.Handle(path, h)
	mux.Handle(path+"/", h)
}

// Init implements interface [surgeon.Initer].
func (s *Server) Init() {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(views.Index()))
	mux.Handle("/auth/", http.StripPrefix("/auth", s.AuthRouter))
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)

	subRoute(mux, "/host", authrouter.RequireAuth, s.HostRouter)

	s.Handler = web.ApplyMiddlewares(mux,
		web.Logger,
		noCache,
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
