package main

import (
	"harmony/internal/server/ioc"
	"net/http"
)

func main() {
	http.ListenAndServe("0.0.0.0:9999", ioc.Server())
}
