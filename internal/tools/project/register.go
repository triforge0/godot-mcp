package project

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "project_info", Description: "Get Godot project metadata.", RPC: "project.info", Level: permission.LevelRead},
		{Name: "project_settings", Description: "Read project settings (subset).", RPC: "project.settings", Level: permission.LevelRead},
		{Name: "project_reload", Description: "Reload the current project.", RPC: "project.reload", Level: permission.LevelWrite},
	})
}
