package server

import (
	"path/filepath"
	"runtime"
)

func projectRoot() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "../..")
}
