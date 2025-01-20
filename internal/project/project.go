package project

import (
	"path/filepath"
	"runtime"
)

func Root() string {
	_, f, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(f), "../..")
}
