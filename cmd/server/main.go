package main

import (
	"context"
	"harmony/cmd/server/ioc"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-cz/devslog"
)

func main() {
	graph := ioc.Root()
	slog.SetDefault(slog.New(devslog.NewHandler(os.Stdout, nil)))

	pump := graph.MessagePump
	server := graph.Server
	err := pump.Start(context.Background())
	if err != nil {
		slog.Error("Error starting pump", "err", err)
		os.Exit(1)
	}

	if err := http.ListenAndServe("0.0.0.0:9999", server); err != nil {
		slog.Error("Error starting http server", "err", err)
		os.Exit(1)
	}
}
