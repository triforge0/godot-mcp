package runtime

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "runtime_run", Description: "Run the project (main scene).", RPC: "runtime.run", Level: permission.LevelExecute},
		{Name: "runtime_stop", Description: "Stop play mode.", RPC: "runtime.stop", Level: permission.LevelExecute},
		{Name: "runtime_pause", Description: "Pause play mode.", RPC: "runtime.pause", Level: permission.LevelExecute},
		{Name: "runtime_resume", Description: "Resume play mode.", RPC: "runtime.resume", Level: permission.LevelExecute},
		{Name: "runtime_status", Description: "Get runtime/play mode status.", RPC: "runtime.status", Level: permission.LevelRead},
		{Name: "runtime_restart", Description: "Restart play mode.", RPC: "runtime.restart", Level: permission.LevelExecute},
	})
}
