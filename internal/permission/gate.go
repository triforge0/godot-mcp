package permission

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

type Gate struct {
	allowDestructive bool
	allowScriptExec  bool
	approver         Approver
	mu               sync.RWMutex
	sessionAlways    map[string]bool // rpc method → always allowed this session
}

func NewGate(approver Approver) *Gate {
	return &Gate{
		allowDestructive: os.Getenv("GODOT_MCP_ALLOW_DESTRUCTIVE") == "1",
		allowScriptExec:  os.Getenv("GODOT_MCP_ALLOW_SCRIPT_EXEC") == "1",
		approver:         approver,
		sessionAlways:    make(map[string]bool),
	}
}

func (g *Gate) Authorize(ctx context.Context, level Level, toolName, rpcMethod string) error {
	if !g.needsPrompt(level, rpcMethod) {
		return nil
	}

	if rpcMethod == "script.execute" && g.allowScriptExec {
		return nil
	}
	if level == LevelDestructive && g.allowDestructive {
		return nil
	}

	if g.isAlwaysAllowed(rpcMethod) {
		return nil
	}

	return g.requestApproval(ctx, level, toolName, rpcMethod)
}

func (g *Gate) needsPrompt(level Level, rpcMethod string) bool {
	if level == LevelDestructive {
		return true
	}
	return rpcMethod == "script.execute"
}

func (g *Gate) isAlwaysAllowed(rpcMethod string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.sessionAlways[rpcMethod]
}

func (g *Gate) setAlwaysAllowed(rpcMethod string) {
	g.mu.Lock()
	g.sessionAlways[rpcMethod] = true
	g.mu.Unlock()
}

func (g *Gate) requestApproval(ctx context.Context, level Level, toolName, rpcMethod string) error {
	if g.approver == nil || !g.approver.Connected() {
		return fmt.Errorf("permission denied: %q requires Godot editor approval (plugin not connected)", toolName)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resp, err := g.approver.RequestPermission(reqCtx, Request{
		Tool:  toolName,
		RPC:   rpcMethod,
		Level: level,
	})
	if err != nil {
		return fmt.Errorf("permission denied: %q — %w", toolName, err)
	}
	if !resp.Approved {
		return fmt.Errorf("permission denied: user declined %q", toolName)
	}
	if resp.Remember == "always" {
		g.setAlwaysAllowed(rpcMethod)
	}
	return nil
}
