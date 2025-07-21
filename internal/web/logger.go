package web

import (
	"harmony/internal/core"
	"log/slog"
	"net/http"
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
		// Don't log request/response cookies
		case "Cookie", "Set-Cookie":
			attrs[i] = slog.Any(k, "...")
		default:
			attrs[i] = slog.Any(k, v)
		}
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
				logHeader(r.Header),
			),
		)
		h.ServeHTTP(rec, r)

		status := rec.Code()
		logLvl := statusCodeToLogLevel(status)
		slog.Log(r.Context(), logLvl, "HTTP Response",
			slog.Group("req",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				logHeader(r.Header),
			),
			slog.Group("res",
				slog.Int("status", status),
				// logHeader(w.Header()),
			),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

var Logger = MiddlewareFunc(Log)
