package console

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "console_logs", Description: "Get recent console log entries.", RPC: "console.logs", Level: permission.LevelRead},
		{Name: "console_clear", Description: "Clear console log buffer.", RPC: "console.clear", Level: permission.LevelWrite},
		{Name: "errors_list", Description: "List recent editor/runtime errors.", RPC: "errors.list", Level: permission.LevelRead},
		{Name: "errors_clear", Description: "Clear error buffer.", RPC: "errors.clear", Level: permission.LevelWrite},
		{Name: "profiler_stats", Description: "Get basic profiler/frame stats.", RPC: "profiler.stats", Level: permission.LevelRead},
	})
}
