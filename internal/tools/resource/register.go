package resource

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "resource_load", Description: "Load a resource by path.", RPC: "resource.load", Level: permission.LevelRead,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "resource_save", Description: "Save a resource.", RPC: "resource.save", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.ResourceCreateParams](r, toolutil.Spec{
		Name: "resource_create", Description: "Create a new resource file.", RPC: "resource.create", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "resource_inspect", Description: "Inspect resource properties.", RPC: "resource.inspect", Level: permission.LevelRead,
	}, nil)
	toolutil.RegisterWithParams[toolutil.ResourceListParams](r, toolutil.Spec{
		Name: "resource_list", Description: "List resources under a path.", RPC: "resource.list", Level: permission.LevelRead,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "resource_delete", Description: "Delete a resource file.", RPC: "resource.delete", Level: permission.LevelDestructive,
	}, nil)
}
