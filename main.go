package main

import (
	"meet-the-locals/views"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := views.Index()

	http.DefaultServeMux.Handle("GET /{$}", templ.Handler(component))
	// http.DefaultServeMux.Handle("GET /static/", http.FileServer(http.FS(static)))
	http.DefaultServeMux.Handle(
		"/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	)

	http.ListenAndServe("0.0.0.0:8081", http.DefaultServeMux)
}
