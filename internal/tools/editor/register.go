package editor

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "editor_selection", Description: "Get current editor selection.", RPC: "editor.selection", Level: permission.LevelRead},
		{Name: "editor_focus", Description: "Focus a node in the editor.", RPC: "editor.focus", Level: permission.LevelWrite},
		{Name: "editor_undo", Description: "Undo last editor action.", RPC: "editor.undo", Level: permission.LevelWrite},
		{Name: "editor_redo", Description: "Redo last editor action.", RPC: "editor.redo", Level: permission.LevelWrite},
		{Name: "editor_refresh", Description: "Refresh editor filesystem.", RPC: "editor.refresh", Level: permission.LevelWrite},
	})
}
