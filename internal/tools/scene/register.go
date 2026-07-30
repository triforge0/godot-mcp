package scene

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "scene_list", Description: "List scenes in the project.", RPC: "scene.list", Level: permission.LevelRead},
		{Name: "scene_current", Description: "Get the currently open scene.", RPC: "scene.current", Level: permission.LevelRead},
		{Name: "scene_open", Description: "Open a scene by res:// path.", RPC: "scene.open", Level: permission.LevelWrite},
		{Name: "scene_save", Description: "Save the current scene.", RPC: "scene.save", Level: permission.LevelWrite},
		{Name: "scene_create", Description: "Create a new empty scene.", RPC: "scene.create", Level: permission.LevelWrite},
		{Name: "scene_reload", Description: "Reload the current scene from disk.", RPC: "scene.reload", Level: permission.LevelWrite},
		{Name: "scene_close", Description: "Close the current scene.", RPC: "scene.close", Level: permission.LevelWrite},
	})
}
