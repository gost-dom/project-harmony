package router

import (
	"context"
	"embed"
	"html/template"
	"io"
	"net/http"

	"harmony/internal/host/router/views"
	serverviews "harmony/internal/web/server/views"

	"github.com/a-h/templ"
)

//go:embed templates/*.*
var fs embed.FS
var templates *template.Template

func init() {
	templates = template.Must(template.ParseFS(fs, "templates/*.tmpl"))
}

type HostRouter struct {
	*http.ServeMux
}

func (r *HostRouter) Index() http.Handler {
	return templ.Handler(views.HostsPage())
}

func (r *HostRouter) RenderPage(page string) http.Handler {
	t := templates.Lookup("new-service-page.tmpl")
	component := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return t.Execute(w, nil)
	})
	l := serverviews.Layout(serverviews.Contents{
		Body: component,
	})
	return templ.Handler(l)
}

// Init implements interface [surgeon.Initer].
func (r *HostRouter) Init() {
	r.ServeMux = http.NewServeMux()
	r.ServeMux.Handle("GET /{$}", r.Index())
	r.ServeMux.Handle("GET /services/new", r.RenderPage(""))
}

func New() *HostRouter {
	r := new(HostRouter)
	r.Init()
	return r
}
