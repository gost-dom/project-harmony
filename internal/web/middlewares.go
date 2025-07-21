package web

import "net/http"

type Middleware func(http.Handler) http.Handler

func JoinMiddlewares(m ...Middleware) Middleware {
	return func(h http.Handler) (res http.Handler) {
		res = h
		for i := len(m) - 1; i >= 0; i-- {
			res = m[i](res)
		}
		return res
	}
}
