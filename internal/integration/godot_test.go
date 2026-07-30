//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/integration"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGodotIntegration(t *testing.T) {
	stack := integration.StartStack(t)
	ctx := context.Background()

	t.Run("plugin_project_info", func(t *testing.T) {
		result, err := stack.Client.CallMap(ctx, "project.info", nil)
		if err != nil {
			t.Fatalf("project.info: %v", err)
		}
		if result["project_name"] != "Godot MCP Demo" {
			t.Fatalf("unexpected project_name: %v", result["project_name"])
		}
		if result["godot_version"] == "" {
			t.Fatal("expected godot_version")
		}
		if result["plugin_version"] == "" {
			t.Fatal("expected plugin_version")
		}
	})

	t.Run("plugin_scene_list", func(t *testing.T) {
		raw, err := stack.Client.Call(ctx, "scene.list", nil, nil)
		if err != nil {
			t.Fatalf("scene.list: %v", err)
		}
		data, err := json.Marshal(raw)
		if err != nil {
			t.Fatal(err)
		}
		if !json.Valid(data) {
			t.Fatalf("invalid json: %s", string(data))
		}
	})

	t.Run("mcp_project_info", func(t *testing.T) {
		res, err := stack.MCPSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "project_info",
		})
		if err != nil {
			t.Fatalf("CallTool project_info: %v", err)
		}
		if res.IsError {
			t.Fatalf("tool error: %v", res.Content)
		}
		text := toolText(res)
		if !strings.Contains(text, "Godot MCP Demo") {
			t.Fatalf("expected project name in result: %s", text)
		}
	})

	t.Run("mcp_scene_open", func(t *testing.T) {
		res, err := stack.MCPSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "scene_open",
			Arguments: map[string]any{
				"path": "res://scenes/main.tscn",
			},
		})
		if err != nil {
			t.Fatalf("CallTool scene_open: %v", err)
		}
		if res.IsError {
			t.Fatalf("tool error: %v", res.Content)
		}
	})
}

func toolText(res *mcpsdk.CallToolResult) string {
	for _, c := range res.Content {
		if tc, ok := c.(*mcpsdk.TextContent); ok {
			return tc.Text
		}
	}
	return ""
}
