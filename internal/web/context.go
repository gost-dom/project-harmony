package web

import (
	"context"
	"net/http"
)

type contextKey string

const (
	CtxKeyCSRFTokenSrc contextKey = "web:csrf-token-src"
	CtxKeyReqID        contextKey = "web:req-id"
)

// ReqValue retrieves a context value from an http Request, and performs a type
// assertion on type T.
func ReqValue[T any](r *http.Request, key any) (res T, ok bool) {
	v := r.Context().Value(key)
	res, ok = v.(T)
	return
}

// SetReqValue takes a reference to a *http.Request variable, and updates the
// variable to point to a new request with a new context extended with the
// key/value v.
func SetReqValue(r **http.Request, key any, v any) {
	ctx := context.WithValue((*r).Context(), key, v)
	*r = (*r).WithContext(ctx)
}
