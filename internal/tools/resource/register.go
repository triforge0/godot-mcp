package resource

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "resource_load", Description: "Load a resource by path.", RPC: "resource.load", Level: permission.LevelRead},
		{Name: "resource_save", Description: "Save a resource.", RPC: "resource.save", Level: permission.LevelWrite},
		{Name: "resource_create", Description: "Create a new resource file.", RPC: "resource.create", Level: permission.LevelWrite},
		{Name: "resource_inspect", Description: "Inspect resource properties.", RPC: "resource.inspect", Level: permission.LevelRead},
		{Name: "resource_list", Description: "List resources under a path.", RPC: "resource.list", Level: permission.LevelRead},
		{Name: "resource_delete", Description: "Delete a resource file.", RPC: "resource.delete", Level: permission.LevelDestructive},
	})
}
