package permission_test

import (
	"context"
	"errors"
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/permission"
)

type mockApprover struct {
	connected bool
	resp      permission.Response
	err       error
	called    bool
	lastReq   permission.Request
}

func (m *mockApprover) Connected() bool { return m.connected }

func (m *mockApprover) RequestPermission(_ context.Context, req permission.Request) (permission.Response, error) {
	m.called = true
	m.lastReq = req
	return m.resp, m.err
}

func TestGateAllowsReadWithoutPrompt(t *testing.T) {
	g := permission.NewGate(nil)
	if err := g.Authorize(context.Background(), permission.LevelRead, "scene_list", "scene.list"); err != nil {
		t.Fatal(err)
	}
}

func TestGateDestructiveRequiresApproval(t *testing.T) {
	mock := &mockApprover{
		connected: true,
		resp:      permission.Response{Approved: true, Remember: "once"},
	}
	g := permission.NewGate(mock)

	if err := g.Authorize(context.Background(), permission.LevelDestructive, "node_delete", "node.delete"); err != nil {
		t.Fatal(err)
	}
	if !mock.called {
		t.Fatal("expected permission request")
	}
	if mock.lastReq.Tool != "node_delete" {
		t.Fatalf("expected node_delete, got %s", mock.lastReq.Tool)
	}
}

func TestGateDestructiveDenied(t *testing.T) {
	mock := &mockApprover{
		connected: true,
		resp:      permission.Response{Approved: false},
	}
	g := permission.NewGate(mock)

	err := g.Authorize(context.Background(), permission.LevelDestructive, "file_delete", "filesystem.delete")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGateAlwaysCached(t *testing.T) {
	mock := &mockApprover{
		connected: true,
		resp:      permission.Response{Approved: true, Remember: "always"},
	}
	g := permission.NewGate(mock)

	if err := g.Authorize(context.Background(), permission.LevelDestructive, "node_delete", "node.delete"); err != nil {
		t.Fatal(err)
	}
	mock.called = false
	if err := g.Authorize(context.Background(), permission.LevelDestructive, "node_delete", "node.delete"); err != nil {
		t.Fatal(err)
	}
	if mock.called {
		t.Fatal("expected cached approval, no second prompt")
	}
}

func TestGateEnvBypassDestructive(t *testing.T) {
	t.Setenv("GODOT_MCP_ALLOW_DESTRUCTIVE", "1")
	mock := &mockApprover{connected: true}
	g := permission.NewGate(mock)

	if err := g.Authorize(context.Background(), permission.LevelDestructive, "node_delete", "node.delete"); err != nil {
		t.Fatal(err)
	}
	if mock.called {
		t.Fatal("env var should bypass dialog")
	}
}

func TestGateScriptExecuteRequiresApproval(t *testing.T) {
	mock := &mockApprover{
		connected: true,
		resp:      permission.Response{Approved: true, Remember: "once"},
	}
	g := permission.NewGate(mock)

	if err := g.Authorize(context.Background(), permission.LevelExecute, "script_execute", "script.execute"); err != nil {
		t.Fatal(err)
	}
	if !mock.called {
		t.Fatal("expected permission request for script.execute")
	}
}

func TestGateNotConnected(t *testing.T) {
	g := permission.NewGate(&mockApprover{connected: false})
	err := g.Authorize(context.Background(), permission.LevelDestructive, "node_delete", "node.delete")
	if err == nil {
		t.Fatal("expected error when plugin disconnected")
	}
}

func TestGateApproverError(t *testing.T) {
	mock := &mockApprover{
		connected: true,
		err:       errors.New("timeout"),
	}
	g := permission.NewGate(mock)
	if err := g.Authorize(context.Background(), permission.LevelDestructive, "node_delete", "node.delete"); err == nil {
		t.Fatal("expected error")
	}
}
