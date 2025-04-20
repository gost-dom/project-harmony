package main

import (
	"harmony/internal/couchdb"
	"harmony/internal/server/ioc"
	"net/http"
)

func main() {
	couchdb.DefaultConnection.StartListener(nil)
	http.ListenAndServe("0.0.0.0:9999", ioc.Server())
}
