package permission

import (
	"context"
	"fmt"
	"os"
)

var blockedActions = map[string]bool{
	"filesystem.delete_project": true,
	"script.execute":            true, // requires explicit execute permission flow
}

type Gate struct {
	allowDestructive bool
	allowScriptExec  bool
}

func NewGate() *Gate {
	return &Gate{
		allowDestructive: os.Getenv("GODOT_MCP_ALLOW_DESTRUCTIVE") == "1",
		allowScriptExec:  os.Getenv("GODOT_MCP_ALLOW_SCRIPT_EXEC") == "1",
	}
}

func (g *Gate) Authorize(_ context.Context, level Level, toolName, rpcMethod string) error {
	if blockedActions[rpcMethod] && !g.allowScriptExec {
		if rpcMethod == "script.execute" {
			return fmt.Errorf("permission denied: %q requires GODOT_MCP_ALLOW_SCRIPT_EXEC=1", toolName)
		}
	}

	switch level {
	case LevelRead:
		return nil
	case LevelWrite, LevelExecute:
		if rpcMethod == "script.execute" && !g.allowScriptExec {
			return fmt.Errorf("permission denied: %q requires GODOT_MCP_ALLOW_SCRIPT_EXEC=1", toolName)
		}
		return nil
	case LevelDestructive:
		if g.allowDestructive {
			return nil
		}
		return fmt.Errorf("permission denied: destructive action %q requires user approval (set GODOT_MCP_ALLOW_DESTRUCTIVE=1 for dev)", toolName)
	default:
		return fmt.Errorf("permission denied: unknown level for %q", toolName)
	}
}
