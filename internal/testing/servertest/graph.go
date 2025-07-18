package servertest

import (
	"harmony/internal/server"
	"harmony/internal/server/ioc"

	"github.com/gost-dom/surgeon"
)

func init() {
	// slog.SetLogLoggerLevel(slog.LevelWarn)
	// logger.SetDefault(slog.Default())
	Graph = ioc.Graph
}

type ServerGraph = *surgeon.Graph[*server.Server]

var Graph *surgeon.Graph[*server.Server]
