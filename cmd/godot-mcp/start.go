package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/godot-mcp/godot-mcp/internal/bridge"
	"github.com/godot-mcp/godot-mcp/internal/config"
	"github.com/godot-mcp/godot-mcp/internal/events"
	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/tools"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the MCP server (stdio) and Godot bridge (WebSocket)",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := config.Load()
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			ws := bridge.NewWebSocket(cfg.BridgeAddr)
			client := godot.NewClient(ws)
			bus := events.NewBus()
			registry := tools.RegisterAll()

			ws.SetEventHandler(func(eventType string, data json.RawMessage) {
				var parsed any
				if len(data) > 0 {
					_ = json.Unmarshal(data, &parsed)
				}
				bus.Publish(eventType, parsed)
			})

			srv := mcp.New(mcp.Options{
				GodotClient: client,
				Registry:    registry,
				EventBus:    bus,
			})

			errCh := make(chan error, 2)
			go func() { errCh <- ws.Start(ctx) }()
			go func() { errCh <- srv.Run(ctx) }()

			select {
			case <-ctx.Done():
				return nil
			case err := <-errCh:
				if err != nil {
					slog.Error("server stopped", "error", err)
					return err
				}
				return nil
			}
		},
	}
}
