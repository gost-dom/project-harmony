package server

import (
	"path/filepath"
	"runtime"
)

func ProjectRoot() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "../..")
}
