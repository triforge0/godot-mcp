//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/integration"
)

// TestAllRPCMethods drives every RPC method the Go tools register against a
// real Godot editor, so a broken adapter (wrong EditorInterface API, a bad
// GDScript signature, etc.) fails here instead of silently shipping.
//
// All mutating work happens inside a disposable scene created by this test,
// never against the checked-in demo scene: opening a different scene while
// one is dirty silently saves it to disk in headless mode (no dialog can
// pop up to ask), so touching examples/demo/scenes/main.tscn here would
// leak test mutations into the repo.
func TestAllRPCMethods(t *testing.T) {
	stack := integration.StartStack(t)
	ctx := context.Background()

	const smokeScenePath = "res://smoke_scene.tscn"
	const skillScenePath = "res://skill_scene.tscn"

	// A dispatch() failure (B.err(...) on the GDScript side) travels back as a
	// JSON-RPC error, which surfaces here as a non-nil Go error — not as an
	// {"error": ...} field in the result map.
	call := func(t *testing.T, method string, params any) map[string]any {
		t.Helper()
		result, err := stack.Client.CallMap(ctx, method, params)
		if err != nil {
			t.Fatalf("%s: %v", method, err)
		}
		return result
	}

	expectErr := func(t *testing.T, method string, params any) {
		t.Helper()
		if _, err := stack.Client.CallMap(ctx, method, params); err == nil {
			t.Fatalf("%s: expected an error, got success", method)
		}
	}

	t.Cleanup(func() {
		_, _ = stack.Client.CallMap(ctx, "filesystem.delete", map[string]string{"path": smokeScenePath})
		_, _ = stack.Client.CallMap(ctx, "filesystem.delete", map[string]string{"path": skillScenePath})
	})

	t.Run("project", func(t *testing.T) {
		call(t, "project.info", nil)
		call(t, "project.settings", nil)
		call(t, "project.reload", nil)
	})

	t.Run("scene", func(t *testing.T) {
		call(t, "scene.list", nil)
		call(t, "scene.open", map[string]string{"path": "res://scenes/main.tscn"})
		call(t, "scene.current", nil)

		// Everything from here on operates on a disposable scene, not main.tscn.
		created := call(t, "scene.create", map[string]string{"name": "SmokeScene", "path": smokeScenePath})
		if created["path"] != smokeScenePath {
			t.Fatalf("scene.create: unexpected path %v", created["path"])
		}
		call(t, "scene.save", nil)
		call(t, "scene.reload", nil)
		expectErr(t, "scene.close", nil) // no EditorInterface API to do this in 4.3
	})

	var rootPath string
	t.Run("node_tree", func(t *testing.T) {
		tree := call(t, "node.list", nil)
		nodes, _ := tree["nodes"].([]any)
		if len(nodes) == 0 {
			t.Fatal("node.list: expected at least the scene root")
		}
		root, _ := nodes[0].(map[string]any)
		rootPath, _ = root["path"].(string)
		if rootPath == "" {
			t.Fatal("node.list: root has no path")
		}

		created := call(t, "node.create", map[string]any{
			"parent_path": rootPath, "type": "Node2D", "name": "SmokeNode",
		})
		nodePath, _ := created["path"].(string)
		if nodePath == "" {
			t.Fatal("node.create: no path returned")
		}

		call(t, "node.get_property", map[string]string{"path": nodePath, "property": "position"})
		call(t, "node.set_property", map[string]any{
			"path": nodePath, "property": "position", "value": map[string]any{"x": 10, "y": 20},
		})
		call(t, "node.children", map[string]string{"path": rootPath})
		call(t, "node.parent", map[string]string{"path": nodePath})

		renamed := call(t, "node.rename", map[string]string{"path": nodePath, "new_name": "SmokeNodeRenamed"})
		renamedPath, _ := renamed["path"].(string)
		if renamedPath == "" {
			t.Fatal("node.rename: no path returned")
		}

		dup := call(t, "node.duplicate", map[string]string{"path": renamedPath})
		dupPath, _ := dup["path"].(string)

		moveTarget := call(t, "node.create", map[string]any{
			"parent_path": rootPath, "type": "Node2D", "name": "MoveTarget",
		})
		moveTargetPath, _ := moveTarget["path"].(string)
		call(t, "node.move", map[string]string{"path": renamedPath, "new_parent_path": moveTargetPath})

		if dupPath != "" {
			call(t, "node.delete", map[string]string{"path": dupPath})
		}
		call(t, "node.delete", map[string]string{"path": moveTargetPath}) // takes the renamed+moved node with it
	})

	t.Run("resource", func(t *testing.T) {
		call(t, "resource.create", map[string]string{"path": "res://smoke_resource.tres", "type": "Resource"})
		call(t, "resource.list", map[string]string{"path": "res://"})
		call(t, "resource.load", map[string]string{"path": "res://smoke_resource.tres"})
		call(t, "resource.inspect", map[string]string{"path": "res://smoke_resource.tres"})
		expectErr(t, "resource.save", map[string]string{"path": "res://smoke_resource.tres"}) // intentionally unimplemented
		call(t, "resource.delete", map[string]string{"path": "res://smoke_resource.tres"})
	})

	t.Run("filesystem", func(t *testing.T) {
		call(t, "filesystem.mkdir", map[string]string{"path": "res://smoke_dir"})
		call(t, "filesystem.create", map[string]string{"path": "res://smoke_dir/file.txt"})
		call(t, "filesystem.write", map[string]string{"path": "res://smoke_dir/file.txt", "content": "hello"})
		read := call(t, "filesystem.read", map[string]string{"path": "res://smoke_dir/file.txt"})
		if read["content"] != "hello" {
			t.Fatalf("filesystem.read: unexpected content %v", read["content"])
		}
		call(t, "filesystem.list", map[string]string{"path": "res://smoke_dir"})
		call(t, "filesystem.delete", map[string]string{"path": "res://smoke_dir/file.txt"})
	})

	t.Run("script", func(t *testing.T) {
		const scriptPath = "res://smoke_script.gd"
		call(t, "script.create", map[string]string{"path": scriptPath, "source": "extends Node2D\n"})
		call(t, "script.read", map[string]string{"path": scriptPath})
		call(t, "script.update", map[string]string{"path": scriptPath, "source": "extends Node2D\n# updated\n"})
		call(t, "script.attach", map[string]string{"path": rootPath, "script_path": scriptPath}) // rootPath == SmokeScene's root
		call(t, "script.detach", map[string]string{"path": rootPath})
		call(t, "script.execute", map[string]string{"source": "1 + 1"})
		call(t, "filesystem.delete", map[string]string{"path": scriptPath})
	})

	t.Run("runtime", func(t *testing.T) {
		call(t, "runtime.status", nil)
		call(t, "runtime.run", nil)
		call(t, "runtime.stop", nil)
		expectErr(t, "runtime.pause", nil) // no EditorInterface API to do this in 4.3
		expectErr(t, "runtime.resume", nil)
	})

	t.Run("screenshot", func(t *testing.T) {
		// Under --headless there's no real rendering server, so viewports have
		// no texture to capture — that's an environment limitation, not a bug.
		for _, method := range []string{"screenshot.capture", "screenshot.viewport"} {
			result, err := stack.Client.CallMap(ctx, method, nil)
			if err != nil {
				t.Logf("%s: capture unavailable (expected under --headless / no GPU): %v", method, err)
				continue
			}
			if result["format"] != "png" {
				t.Fatalf("%s: unexpected format %v", method, result["format"])
			}
		}
	})

	t.Run("editor", func(t *testing.T) {
		call(t, "editor.selection", nil)
		call(t, "editor.focus", map[string]string{"path": rootPath})
		call(t, "editor.undo", nil)
		call(t, "editor.redo", nil)
		call(t, "editor.refresh", nil)
	})

	t.Run("console", func(t *testing.T) {
		call(t, "console.logs", nil)
		call(t, "errors.list", nil)
		call(t, "profiler.stats", nil)
		call(t, "console.clear", nil)
		call(t, "errors.clear", nil)
	})

	t.Run("reflection", func(t *testing.T) {
		call(t, "object.inspect", map[string]string{"path": rootPath})
		call(t, "class.inspect", map[string]string{"class": "Node2D"})
		call(t, "property.list", map[string]string{"path": rootPath})
		call(t, "method.list", map[string]string{"path": rootPath})
	})

	t.Run("skills", func(t *testing.T) {
		// parent_path left empty: skill.create_player defaults to whatever
		// scene is currently edited, so this doesn't depend on rootPath
		// pointing at the right scene after create_scene switches tabs.
		call(t, "skill.create_scene", map[string]string{"name": "SkillScene", "path": skillScenePath})
		call(t, "skill.create_player", nil)
		call(t, "skill.optimize_project", nil)
		call(t, "skill.analyze_error", map[string]int{"limit": 5})
	})
}
