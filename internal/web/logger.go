package web

import (
	"fmt"
	"harmony/internal/core"
	"harmony/internal/infrastructure/log"
	"net/http"
	"slices"
	"strings"
	"time"
)

func statusCodeToLogLevel(code int) log.Level {
	if code >= 500 {
		return log.LevelError
	}
	if code >= 400 {
		return log.LevelWarn
	}
	return log.LevelInfo
}

// logHeader creates an [log.Attr] representing HTTP request or response
// headers. Cookies values are hidden, but cookie names and options are kept to
// debug malfunctioning cookies.
func logHeader(h http.Header) log.Attr {
	attrs := make([]any, len(h))
	i := 0
	for k, v := range h {
		switch k {
		case "Cookie", "Set-Cookie":
			v = slices.Clone(v)
			for j := range v {
				components := strings.Split(v[j], ";")
				v[j] = "******"
				if len(components) > 0 {
					parts := strings.Split(components[0], "=")
					if len(parts) > 0 {
						components[0] = fmt.Sprintf("%s=******", parts[0])
						v[j] = strings.Join(components, ";")
					}
				}
			}
		}
		attrs[i] = log.Any(k, v)
		i++
	}
	return log.Group("header", attrs...)
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.ContextWith(&r, "reqID", core.NewID())

		rec := &StatusRecorder{ResponseWriter: w}
		start := time.Now()

		log.Info(r.Context(), "HTTP Request",
			log.Group("req",
				log.String("method", r.Method),
				log.String("path", r.URL.Path),
			),
		)
		log.Debug(r.Context(), "HTTP Request headers", logHeader(r.Header))
		h.ServeHTTP(rec, r)

		status := rec.Code()
		logLvl := statusCodeToLogLevel(status)
		log.Log(r.Context(), logLvl, "HTTP Response",
			log.Group("req",
				log.String("method", r.Method),
				log.String("path", r.URL.Path),
			),
			log.Group("res",
				log.Int("status", status),
			),
			log.Duration("duration", time.Since(start)),
		)
		log.Debug(r.Context(), "HTTP Response headers", logHeader(w.Header()))
	})
}

var Logger = MiddlewareFunc(Log)
