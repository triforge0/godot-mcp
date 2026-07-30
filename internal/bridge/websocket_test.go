package bridge_test

import (
	"context"
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/bridge"
)

func TestWebSocketNotConnected(t *testing.T) {
	ws := bridge.NewWebSocket("127.0.0.1:1")
	_, err := ws.Call(context.Background(), "project.info", nil)
	if err != bridge.ErrNotConnected {
		t.Fatalf("expected ErrNotConnected, got %v", err)
	}
}
