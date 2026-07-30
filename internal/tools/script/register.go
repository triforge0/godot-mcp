package script

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "script_read", Description: "Read a GDScript file.", RPC: "script.read", Level: permission.LevelRead,
	}, nil)
	toolutil.RegisterWithParams[toolutil.ScriptSourceParams](r, toolutil.Spec{
		Name: "script_create", Description: "Create a new GDScript file.", RPC: "script.create", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.ScriptSourceParams](r, toolutil.Spec{
		Name: "script_update", Description: "Update a GDScript file.", RPC: "script.update", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.ScriptAttachParams](r, toolutil.Spec{
		Name: "script_attach", Description: "Attach a script to a node.", RPC: "script.attach", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "script_detach", Description: "Detach script from a node.", RPC: "script.detach", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.ScriptSourceParams](r, toolutil.Spec{
		Name: "script_execute", Description: "Execute GDScript snippet in editor context.", RPC: "script.execute", Level: permission.LevelExecute,
	}, nil)
}
