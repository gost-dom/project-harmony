package main

import (
	"harmony/internal/server/ioc"
	"net/http"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func main() {
	http.ListenAndServe("0.0.0.0:8081", ioc.Server())
}
