package servertest

import (
	"harmony/internal/server"
	"harmony/internal/server/ioc"

	"github.com/gost-dom/surgeon"
)

func init() {
	// slog.SetLogLoggerLevel(slog.LevelWarn)
	// logger.SetDefault(slog.Default())
	graph = ioc.Graph
}

var graph *surgeon.Graph[*server.Server]
