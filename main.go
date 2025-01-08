package main

import "net/http"

func main() {
	server := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("<html><body><h1>Foo</h1></body></html>"))
	})
	http.ListenAndServe("0.0.0.0:3000", server)
}
