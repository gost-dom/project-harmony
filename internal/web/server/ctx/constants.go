package serverctx

import (
	"context"
	"net/http"
)

type ContextKey string

const (
	ServerCSRFTokenSrc ContextKey = "server:csrf:token-source"
	ServerReqID        ContextKey = "server:req-id"
)

func ReqValue[T any](r *http.Request, key any) (res T, ok bool) {
	v := r.Context().Value(key)
	res, ok = v.(T)
	return
}

func SetReqValue(r **http.Request, key any, v any) {
	ctx := context.WithValue((*r).Context(), key, v)
	*r = (*r).WithContext(ctx)
}
