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

// Init implements interface [surgeon.Initer].
func (r *HostRouter) Init() {
	r.ServeMux = http.NewServeMux()
	r.ServeMux.Handle("GET /", r.Index())
}

func New() *HostRouter {
	r := new(HostRouter)
	r.Init()
	return r
}
