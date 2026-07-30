package screenshot

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterAll(r, []toolutil.Spec{
		{Name: "screenshot_capture", Description: "Capture editor or game viewport as PNG (base64).", RPC: "screenshot.capture", Level: permission.LevelRead},
		{Name: "screenshot_viewport", Description: "Capture a specific viewport.", RPC: "screenshot.viewport", Level: permission.LevelRead},
	})
}
