// Package godotmcp is the public SDK for registering custom Godot MCP tools.
package godotmcp

import (
	"context"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

type (
	Registry       = mcp.Registry
	ToolDefinition = mcp.ToolDefinition
	HandlerFunc    = mcp.HandlerFunc
	Client         = godot.Client
	Level          = permission.Level
	Spec           = toolutil.Spec
)

const (
	LevelRead        = permission.LevelRead
	LevelWrite       = permission.LevelWrite
	LevelDestructive = permission.LevelDestructive
	LevelExecute     = permission.LevelExecute
)

func NewRegistry() *Registry { return mcp.NewRegistry() }

func RegisterTool(r *Registry, def ToolDefinition) { r.Register(def) }

func RegisterNoParams(r *Registry, spec Spec) { toolutil.RegisterNoParams(r, spec) }

func RegisterWithParams[T any](r *Registry, spec Spec, mapParams func(T) any) {
	toolutil.RegisterWithParams(r, spec, mapParams)
}

type HandlerContext = context.Context
