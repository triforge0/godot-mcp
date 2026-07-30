package tools

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/tools/console"
	"github.com/godot-mcp/godot-mcp/internal/tools/editor"
	"github.com/godot-mcp/godot-mcp/internal/tools/filesystem"
	"github.com/godot-mcp/godot-mcp/internal/tools/node"
	"github.com/godot-mcp/godot-mcp/internal/tools/project"
	"github.com/godot-mcp/godot-mcp/internal/tools/reflection"
	"github.com/godot-mcp/godot-mcp/internal/tools/resource"
	"github.com/godot-mcp/godot-mcp/internal/tools/runtime"
	"github.com/godot-mcp/godot-mcp/internal/tools/scene"
	"github.com/godot-mcp/godot-mcp/internal/tools/screenshot"
	"github.com/godot-mcp/godot-mcp/internal/tools/script"
	"github.com/godot-mcp/godot-mcp/internal/tools/skills"
)

func RegisterAll() *mcp.Registry {
	r := mcp.NewRegistry()
	project.Register(r)
	scene.Register(r)
	node.Register(r)
	resource.Register(r)
	filesystem.Register(r)
	script.Register(r)
	runtime.Register(r)
	screenshot.Register(r)
	editor.Register(r)
	console.Register(r)
	reflection.Register(r)
	skills.Register(r)
	return r
}

func Count() int {
	return len(RegisterAll().All())
}
