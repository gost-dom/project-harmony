package main

import (
	"fmt"
	"harmony/internal/core/corerepo"
	"harmony/internal/messaging"
	mioc "harmony/internal/messaging/ioc"
	"harmony/internal/server/ioc"
	"log/slog"
	"net/http"
)

func main() {
	fmt.Println("Starting server")
	pump := messaging.MessagePump{
		MessageSource:         corerepo.DefaultMessageSource,
		DomainEventRepository: corerepo.DefaultDomainEventRepo,
		Handler:               mioc.Handler(),
	}
	err := pump.Start(nil)
	if err != nil {
		slog.Error("Error starting pump", "err", err)
	}
	http.ListenAndServe("0.0.0.0:9999", ioc.Server())
}
