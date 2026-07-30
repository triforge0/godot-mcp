package e2e_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/godot-mcp/godot-mcp/internal/bridge"
	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
	"github.com/gorilla/websocket"
)

func mockGodotPlugin(t *testing.T, wsURL string) {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial mock plugin: %v", err)
	}
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var req bridge.Request
			if json.Unmarshal(msg, &req) != nil {
				continue
			}
			var result any
			switch req.Method {
			case "project.info":
				result = map[string]string{
					"plugin": "godot-mcp", "plugin_version": "1.0.0",
					"godot_version": "4.3.0", "project_name": "Test",
				}
			case "scene.open":
				var p map[string]string
				_ = json.Unmarshal(req.Params, &p)
				if p["path"] == "" {
					writeError(conn, req.ID, "path required")
					continue
				}
				result = map[string]any{"opened": true, "path": p["path"]}
			case "node.create":
				var p map[string]string
				_ = json.Unmarshal(req.Params, &p)
				result = map[string]any{
					"class": p["type"], "name": p["name"],
					"path": p["parent_path"] + "/" + p["name"],
				}
			case "node.set_property":
				var p map[string]any
				_ = json.Unmarshal(req.Params, &p)
				result = map[string]any{"path": p["path"], "property": p["property"], "value": p["value"]}
			case "runtime.run":
				result = map[string]any{"running": true}
			case "runtime.stop":
				result = map[string]any{"stopped": true}
			case "screenshot.capture":
				result = map[string]any{"format": "png", "data": "aGVsbG8=", "width": 1, "height": 1}
			default:
				result = map[string]bool{"ok": true}
			}
			writeResult(conn, req.ID, result)
		}
	}()
}

func writeResult(conn *websocket.Conn, id string, result any) {
	raw, _ := json.Marshal(result)
	resp, _ := json.Marshal(bridge.Response{JSONRPC: bridge.JSONRPCVersion, ID: id, Result: raw})
	_ = conn.WriteMessage(websocket.TextMessage, resp)
}

func writeError(conn *websocket.Conn, id, msg string) {
	resp, _ := json.Marshal(bridge.Response{
		JSONRPC: bridge.JSONRPCVersion, ID: id,
		Error: &bridge.Error{Code: -32000, Message: msg},
	})
	_ = conn.WriteMessage(websocket.TextMessage, resp)
}

func TestE2E_TopTools(t *testing.T) {
	const addr = "127.0.0.1:19506"
	ws := bridge.NewWebSocket(addr)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go func() { _ = ws.Start(ctx) }()
	time.Sleep(50 * time.Millisecond)

	mockGodotPlugin(t, "ws://"+addr+"/ws")
	deadline := time.Now().Add(2 * time.Second)
	for !ws.Connected() && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if !ws.Connected() {
		t.Fatal("mock plugin did not connect")
	}

	client := godot.NewClient(ws)
	registry := mcp.NewRegistry()

	toolutil.RegisterWithParams[toolutil.PathParams](registry, toolutil.Spec{
		Name: "scene_open", RPC: "scene.open", Level: permission.LevelWrite, Description: "open",
	}, nil)
	toolutil.RegisterWithParams[toolutil.NodeCreateParams](registry, toolutil.Spec{
		Name: "node_create", RPC: "node.create", Level: permission.LevelWrite, Description: "create",
	}, nil)
	toolutil.RegisterWithParams[toolutil.PropertySetParams](registry, toolutil.Spec{
		Name: "node_set_property", RPC: "node.set_property", Level: permission.LevelWrite, Description: "set",
	}, nil)
	toolutil.RegisterNoParams(registry, toolutil.Spec{
		Name: "runtime_run", RPC: "runtime.run", Level: permission.LevelExecute, Description: "run",
	})
	toolutil.RegisterNoParams(registry, toolutil.Spec{
		Name: "runtime_stop", RPC: "runtime.stop", Level: permission.LevelExecute, Description: "stop",
	})
	toolutil.RegisterNoParams(registry, toolutil.Spec{
		Name: "screenshot_capture", RPC: "screenshot.capture", Level: permission.LevelRead, Description: "shot",
	})

	cases := []struct {
		name   string
		params any
	}{
		{"scene_open", toolutil.PathParams{Path: "res://main.tscn"}},
		{"node_create", toolutil.NodeCreateParams{ParentPath: "Root", Type: "Sprite2D", Name: "Player"}},
		{"node_set_property", toolutil.PropertySetParams{Path: "Root/Player", Property: "position", Value: map[string]any{"x": 10, "y": 20}}},
		{"runtime_run", nil},
		{"runtime_stop", nil},
		{"screenshot_capture", nil},
	}

	for _, tc := range cases {
		def, ok := registry.ByName(tc.name)
		if !ok {
			t.Fatalf("missing tool %s", tc.name)
		}
		if _, err := def.Handler(context.Background(), client, tc.params); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
	}

	info := ws.PluginInfo()
	if info == nil || info.PluginVersion == "" {
		t.Log("plugin info not yet cached (async fetch)")
	}
}
