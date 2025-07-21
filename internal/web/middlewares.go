package web

import "net/http"

type Middleware interface {
	Apply(http.Handler) http.Handler
}

type MiddlewareFunc func(http.Handler) http.Handler

func (f MiddlewareFunc) Apply(h http.Handler) http.Handler { return f(h) }

func ApplyMiddlewares(h http.Handler, m ...Middleware) (res http.Handler) {
	res = h
	for i := len(m) - 1; i >= 0; i-- {
		res = m[i].Apply(res)
	}
	return res
}
