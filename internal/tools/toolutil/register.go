package toolutil

import (
	"context"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
)

type Spec struct {
	Name        string
	Description string
	RPC         string
	Level       permission.Level
}

func RegisterRPC(registry *mcp.Registry, spec Spec) {
	registry.Register(mcp.ToolDefinition{
		Name:        spec.Name,
		Description: spec.Description,
		RPCMethod:   spec.RPC,
		Permission:  spec.Level,
		Handler: func(ctx context.Context, client *godot.Client, params any) (any, error) {
			return client.Call(ctx, spec.RPC, params, nil)
		},
	})
}

func RegisterAll(registry *mcp.Registry, specs []Spec) {
	for _, spec := range specs {
		RegisterRPC(registry, spec)
	}
}
