package log

import (
	"context"
	"log/slog"
)

type contextKey string

const ctxKeyLogger contextKey = "infra:log:logger"

func logger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKeyLogger).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func Log(ctx context.Context, level Level, msg string, args ...any) {
	l := logger(ctx)
	l.Log(ctx, level, msg, args...)
}

// Contextable is an interface for a type that can hold a context, and for which
// you can create a clone with a new child context.
type Contextable[T any] interface {
	Context() context.Context
	WithContext(context.Context) T
}

// ContextWith update a reference to a context-bearing value, setting a new
// value with to contain an updated logger always outputting attributes from
// args.
//
// E.g., *http.Request is a compatible context-bearing type.
func ContextWith[T Contextable[T]](cp *T, args ...any) {
	c := *cp
	*cp = c.WithContext(With(c.Context(), args...))
}

func With(ctx context.Context, args ...any) context.Context {
	l := logger(ctx)
	return context.WithValue(ctx, ctxKeyLogger, l.With(args...))
}

func WithGroup(ctx context.Context, name string) context.Context {
	l := logger(ctx)
	return context.WithValue(ctx, ctxKeyLogger, l.WithGroup(name))
}

func Debug(ctx context.Context, msg string, args ...any) {
	Log(ctx, LevelDebug, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	Log(ctx, LevelInfo, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	Log(ctx, LevelWarn, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	Log(ctx, LevelError, msg, args...)
}

// Re-exporting names from slog, makes it easier to grep for usages in slog in
// the system to find what needs to be replaced; or create a rule to not use
// slog.

type Level = slog.Level
type Attr = slog.Attr

var Group = slog.Group
var Int = slog.Int
var Duration = slog.Duration
var String = slog.String
var Any = slog.Any

var LevelError = slog.LevelError
var LevelInfo = slog.LevelInfo
var LevelWarn = slog.LevelWarn
var LevelDebug = slog.LevelDebug
