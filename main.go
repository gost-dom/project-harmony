package main

import (
	"harmony/views"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	// server := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	// 	res.Write([]byte("<html><body><h1>Foo bar!</h1></body></html>"))
	// })

	component := views.Index()

	server := templ.Handler(component)

	http.ListenAndServe("0.0.0.0:8081", server)
}
