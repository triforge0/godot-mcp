package toolutil_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

type mockCaller struct {
	lastMethod string
	lastParams any
}

func (m *mockCaller) Call(_ context.Context, method string, params any) (json.RawMessage, error) {
	m.lastMethod = method
	m.lastParams = params
	return json.RawMessage(`{"ok":true}`), nil
}

func (m *mockCaller) Connected() bool { return true }

func TestRegisterWithParamsPassesArgs(t *testing.T) {
	caller := &mockCaller{}
	client := godot.NewClient(caller)

	registry := mcp.NewRegistry()
	toolutil.RegisterWithParams[toolutil.PathParams](registry, toolutil.Spec{
		Name: "scene_open", RPC: "scene.open", Level: permission.LevelWrite,
		Description: "open scene",
	}, nil)

	def, ok := registry.ByName("scene_open")
	if !ok {
		t.Fatal("tool not registered")
	}

	_, err := def.Handler(context.Background(), client, toolutil.PathParams{Path: "res://main.tscn"})
	if err != nil {
		t.Fatal(err)
	}
	if caller.lastMethod != "scene.open" {
		t.Fatalf("expected scene.open, got %s", caller.lastMethod)
	}
	params, ok := caller.lastParams.(toolutil.PathParams)
	if !ok || params.Path != "res://main.tscn" {
		t.Fatalf("expected path param, got %#v", caller.lastParams)
	}
}

func TestRegisterNoParamsSendsNil(t *testing.T) {
	caller := &mockCaller{}
	client := godot.NewClient(caller)

	registry := mcp.NewRegistry()
	toolutil.RegisterNoParams(registry, toolutil.Spec{
		Name: "runtime_run", RPC: "runtime.run", Level: permission.LevelExecute,
		Description: "run",
	})

	def, _ := registry.ByName("runtime_run")
	_, err := def.Handler(context.Background(), client, nil)
	if err != nil {
		t.Fatal(err)
	}
	if caller.lastParams != nil {
		t.Fatalf("expected nil params, got %#v", caller.lastParams)
	}
}
