package mcp

import (
	"context"
	"log/slog"

	"github.com/godot-mcp/godot-mcp/internal/events"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Notifier forwards EventBus events to all MCP sessions as logging notifications.
type Notifier struct {
	server *mcpsdk.Server
}

func NewNotifier(server *mcpsdk.Server, bus *events.Bus) *Notifier {
	n := &Notifier{server: server}
	bus.Subscribe("*", n.handleEvent)
	return n
}

func (n *Notifier) handleEvent(event events.Event) {
	payload := map[string]any{
		"type": event.Type,
		"data": event.Data,
		"ts":   event.Timestamp,
	}

	ctx := context.Background()
	for session := range n.server.Sessions() {
		if err := session.Log(ctx, &mcpsdk.LoggingMessageParams{
			Level:  "info",
			Logger: "godot-mcp",
			Data:   payload,
		}); err != nil {
			slog.Debug("mcp notification failed", "error", err)
		}
	}
}

// BridgeEventTypes are events emitted by the Godot plugin.
var BridgeEventTypes = []string{
	"scene.changed",
	"node.selected",
	"runtime.started",
	"runtime.stopped",
	"runtime.error",
}
