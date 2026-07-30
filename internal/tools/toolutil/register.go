package toolutil

import (
	"context"
	"fmt"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Spec struct {
	Name        string
	Description string
	RPC         string
	Level       permission.Level
}

// RegisterNoParams registers a tool with no MCP input parameters.
func RegisterNoParams(registry *mcp.Registry, spec Spec) {
	registry.Register(mcp.ToolDefinition{
		Name:        spec.Name,
		Description: spec.Description,
		RPCMethod:   spec.RPC,
		Permission:  spec.Level,
		Handler: func(ctx context.Context, client *godot.Client, _ any) (any, error) {
			return client.Call(ctx, spec.RPC, nil, nil)
		},
		MCPRegister: mcp.RegisterNoArgs,
	})
}

// RegisterWithParams registers a tool with typed MCP input parameters.
func RegisterWithParams[T any](registry *mcp.Registry, spec Spec, mapParams func(T) any) {
	if mapParams == nil {
		mapParams = func(t T) any { return t }
	}
	registry.Register(mcp.ToolDefinition{
		Name:        spec.Name,
		Description: spec.Description,
		RPCMethod:   spec.RPC,
		Permission:  spec.Level,
		Handler: func(ctx context.Context, client *godot.Client, params any) (any, error) {
			p, ok := params.(T)
			if !ok {
				return nil, fmt.Errorf("invalid params for %s", spec.Name)
			}
			return client.Call(ctx, spec.RPC, mapParams(p), nil)
		},
		MCPRegister: func(server *mcpsdk.Server, def mcp.ToolDefinition, dispatch mcp.DispatchFunc) {
			mcp.RegisterToolWithArgs(server, def, dispatch, func(a T) any {
				return mapParams(a)
			})
		},
	})
}

// RegisterAll registers multiple no-parameter tools.
func RegisterAll(registry *mcp.Registry, specs []Spec) {
	for _, spec := range specs {
		RegisterNoParams(registry, spec)
	}
}
