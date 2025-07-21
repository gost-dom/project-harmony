package main

import (
	"context"
	"fmt"
	"harmony/internal/core/corerepo"
	"harmony/internal/messaging"
	mioc "harmony/internal/messaging/ioc"
	"harmony/internal/web/server/ioc"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-cz/devslog"
)

func main() {
	slog.SetDefault(slog.New(devslog.NewHandler(os.Stdout, nil)))

	fmt.Println("Starting server")
	pump := messaging.MessagePump{
		MessageSource:         corerepo.DefaultMessageSource,
		DomainEventRepository: corerepo.DefaultDomainEventRepo,
		Handler:               mioc.Handler(),
	}
	err := pump.Start(context.Background())
	if err != nil {
		slog.Error("Error starting pump", "err", err)
	}

	http.ListenAndServe("0.0.0.0:9999", ioc.Server())
}
