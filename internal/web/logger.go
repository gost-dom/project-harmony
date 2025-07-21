package web

import (
	"fmt"
	"harmony/internal/core"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func statusCodeToLogLevel(code int) slog.Level {
	if code >= 500 {
		return slog.LevelError
	}
	if code >= 400 {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func logHeader(h http.Header) slog.Attr {
	attrs := make([]any, len(h))
	i := 0
	for k, v := range h {
		switch k {
		// Hide cookie values
		case "Cookie", "Set-Cookie":
			u := v
			v = make([]string, len(v))
			for j := range v {
				v[j] = "******"
				components := strings.Split(u[j], ";")
				if len(components) > 0 {
					parts := strings.Split(components[0], "=")
					if len(parts) > 0 {
						components[0] = fmt.Sprintf("%s=******", parts[0])
						v[j] = strings.Join(components, ";")
					}
				}
			}
		}
		attrs[i] = slog.Any(k, v)
		i++
	}
	return slog.Group("header", attrs...)
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetReqValue(&r, CtxKeyReqID, core.NewID())
		rec := &StatusRecorder{ResponseWriter: w}
		start := time.Now()

		slog.InfoContext(r.Context(), "HTTP Request",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			),
		)
		slog.Log(r.Context(), slog.LevelDebug, "HTTP Request headers", logHeader(r.Header))
		h.ServeHTTP(rec, r)

		status := rec.Code()
		logLvl := statusCodeToLogLevel(status)
		slog.Log(r.Context(), logLvl, "HTTP Response",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			),
			slog.Group("res",
				slog.Int("status", status),
			),
			slog.Duration("duration", time.Since(start)),
		)
		slog.Log(r.Context(), slog.LevelDebug, "HTTP Response headers", logHeader(w.Header()))
	})
}

var Logger = MiddlewareFunc(Log)
