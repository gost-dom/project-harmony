package hostrouter

import (
	"harmony/internal/host/hostrouter/views"
	"net/http"

	"github.com/a-h/templ"
)

type HostRouter struct {
	*http.ServeMux
}

func (r *HostRouter) Index() http.Handler {
	return templ.Handler(views.HostsPage())
}

func (r *HostRouter) Init() {
	r.ServeMux = http.NewServeMux()
	r.ServeMux.Handle("GET /", templ.Handler(views.HostsPage()))
}

func NewHostRouter() *HostRouter {
	r := new(HostRouter)
	r.Init()
	return r
}
