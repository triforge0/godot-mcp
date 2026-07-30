package reflection

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "object_inspect", Description: "Inspect any Godot object by path.", RPC: "object.inspect", Level: permission.LevelRead},
		{Name: "class_inspect", Description: "Inspect a Godot class via ClassDB.", RPC: "class.inspect", Level: permission.LevelRead},
		{Name: "property_list", Description: "List properties on an object.", RPC: "property.list", Level: permission.LevelRead},
		{Name: "method_list", Description: "List methods on an object.", RPC: "method.list", Level: permission.LevelRead},
	})
}
