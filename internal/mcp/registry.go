package mcp

import (
	"context"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type HandlerFunc func(ctx context.Context, client *godot.Client, params any) (any, error)

type MCPRegistrar func(server *mcpsdk.Server, def ToolDefinition, dispatch DispatchFunc)

type ToolDefinition struct {
	Name        string
	Description string
	RPCMethod   string
	Permission  permission.Level
	Handler     HandlerFunc
	MCPRegister MCPRegistrar
}

type Registry struct {
	tools []ToolDefinition
}

func NewRegistry() *Registry {
	return &Registry{tools: make([]ToolDefinition, 0, 64)}
}

func (r *Registry) Register(def ToolDefinition) {
	if def.Handler == nil {
		panic("tool " + def.Name + " has no handler")
	}
	if def.MCPRegister == nil {
		def.MCPRegister = RegisterNoArgs
	}
	r.tools = append(r.tools, def)
}

func (r *Registry) All() []ToolDefinition {
	return append([]ToolDefinition(nil), r.tools...)
}

func (r *Registry) ByName(name string) (ToolDefinition, bool) {
	for _, def := range r.tools {
		if def.Name == name {
			return def, true
		}
	}
	return ToolDefinition{}, false
}

func (r *Registry) Len() int { return len(r.tools) }
