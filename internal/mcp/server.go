package mcp

import (
	"context"
	"log/slog"

	"github.com/godot-mcp/godot-mcp/internal/events"
	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/version"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcp *mcpsdk.Server
	bus *events.Bus
}

type Options struct {
	GodotClient *godot.Client
	Registry    *Registry
	EventBus    *events.Bus
}

func New(opts Options) *Server {
	bus := opts.EventBus
	if bus == nil {
		bus = events.NewBus()
	}

	gate := permission.NewGate(opts.GodotClient)
	dispatcher := NewDispatcher(opts.GodotClient, gate, bus)

	mcpServer := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    version.Name,
		Version: version.Version,
	}, nil)

	RegisterTools(mcpServer, opts.Registry, dispatcher)
	NewNotifier(mcpServer, bus)

	slog.Info("mcp server initialized", "tools", opts.Registry.Len(), "version", version.Version)

	return &Server{mcp: mcpServer, bus: bus}
}

func (s *Server) EventBus() *events.Bus { return s.bus }

func (s *Server) MCPServer() *mcpsdk.Server { return s.mcp }

func (s *Server) Run(ctx context.Context) error {
	slog.Info("mcp server starting", "transport", "stdio")
	return s.mcp.Run(ctx, &mcpsdk.StdioTransport{})
}
