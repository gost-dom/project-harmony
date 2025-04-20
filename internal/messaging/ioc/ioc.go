package ioc

import (
	"harmony/internal/messaging"

	"github.com/gost-dom/surgeon"
)

var Graph *surgeon.Graph[*messaging.MessageHandler]

func init() {
	Graph = surgeon.BuildGraph(&messaging.MessageHandler{})
}

func Handler() *messaging.MessageHandler { return Graph.Instance() }
