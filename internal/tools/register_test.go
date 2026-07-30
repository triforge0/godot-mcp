package tools_test

import (
	"context"
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/tools"
)

func TestToolCount(t *testing.T) {
	n := tools.RegisterAll().Len()
	if n < 40 {
		t.Fatalf("expected at least 40 tools, got %d", n)
	}
}

func TestParameterizedToolsRejectNilParams(t *testing.T) {
	r := tools.RegisterAll()
	client := godot.NewClient(nil)
	ctx := context.Background()

	paramTools := []string{
		"scene_open", "node_create", "node_set_property",
		"file_read", "script_attach", "editor_focus",
	}
	for _, name := range paramTools {
		def, ok := r.ByName(name)
		if !ok {
			t.Fatalf("missing tool %s", name)
		}
		if def.MCPRegister == nil {
			t.Fatalf("tool %s has no MCP registrar", name)
		}
		_, err := def.Handler(ctx, client, nil)
		if err == nil {
			t.Fatalf("tool %s should reject nil params", name)
		}
	}
}

func TestNoParamToolsUseDefaultRegistrar(t *testing.T) {
	r := tools.RegisterAll()
	def, ok := r.ByName("project_info")
	if !ok {
		t.Fatal("missing project_info")
	}
	if def.MCPRegister == nil {
		t.Fatal("expected default MCP registrar")
	}
}
