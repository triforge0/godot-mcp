package tools_test

import (
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/tools"
)

func TestToolCount(t *testing.T) {
	n := tools.RegisterAll().Len()
	if n < 40 {
		t.Fatalf("expected at least 40 tools, got %d", n)
	}
}
