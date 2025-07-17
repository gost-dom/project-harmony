package htest

import (
	"log/slog"
	"testing"

	"github.com/gost-dom/browser/testing/gosttest"
)

// UseTestLogger installs a default slog logger and use Cleanup to restore the
// original logger.
//
// Because UseTestLogger affects the whole process, it cannot be used in
// parallel tests or tests with parallel ancestors.
func UseTestLogger(t testing.TB) {
	orig := slog.Default()
	t.Cleanup(func() { slog.SetDefault(orig) })

	l := gosttest.NewTestingLogger(t)
	slog.SetDefault(l)
}
