package node

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterNoParams(r, toolutil.Spec{
		Name: "node_list", Description: "List nodes in the current scene tree.", RPC: "node.list", Level: permission.LevelRead,
	})

	toolutil.RegisterWithParams[toolutil.NodeCreateParams](r, toolutil.Spec{
		Name: "node_create", Description: "Create a node (any ClassDB type).", RPC: "node.create", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "node_delete", Description: "Delete a node.", RPC: "node.delete", Level: permission.LevelDestructive,
	}, nil)

	toolutil.RegisterWithParams[toolutil.NodeRenameParams](r, toolutil.Spec{
		Name: "node_rename", Description: "Rename a node.", RPC: "node.rename", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterWithParams[toolutil.NodeMoveParams](r, toolutil.Spec{
		Name: "node_move", Description: "Reparent a node.", RPC: "node.move", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "node_duplicate", Description: "Duplicate a node.", RPC: "node.duplicate", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterWithParams[toolutil.PropertyGetParams](r, toolutil.Spec{
		Name: "node_get_property", Description: "Get a node property via reflection.", RPC: "node.get_property", Level: permission.LevelRead,
	}, nil)

	toolutil.RegisterWithParams[toolutil.PropertySetParams](r, toolutil.Spec{
		Name: "node_set_property", Description: "Set a node property via reflection.", RPC: "node.set_property", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "node_children", Description: "List child nodes.", RPC: "node.children", Level: permission.LevelRead,
	}, nil)

	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "node_parent", Description: "Get parent node info.", RPC: "node.parent", Level: permission.LevelRead,
	}, nil)
}
