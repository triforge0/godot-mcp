# Godot MCP Demo

Minimal Godot 4 project for testing [godot-mcp](../../README.md).

## Quick start

1. Open this folder in Godot 4.3+ (`examples/demo/`).
2. Enable **Project → Project Settings → Plugins → Godot MCP**.
3. Build and configure the MCP server (see [GETTING_STARTED.md](../../docs/GETTING_STARTED.md)).
4. Run `make doctor` from the repo root with Godot open.

## Try these tools

| Tool | Example prompt |
|------|----------------|
| `project_info` | "Show Godot project metadata" |
| `scene_list` | "List all scenes in the project" |
| `scene_open` | "Open `res://scenes/main.tscn`" |
| `node_list` | "List nodes in the current scene" |
| `node_set_property` | "Set the Label text to Hello MCP" |
| `screenshot_capture` | "Capture the editor viewport" |
| `runtime_run` | "Start play mode" |

Destructive actions (`node_delete`, `file_delete`, etc.) and `script_execute` show a **permission dialog** in the Godot editor.

## Plugin path

The addon is symlinked from `../../plugin/addons/godot_mcp`. If the symlink is missing:

```bash
ln -sfn ../../../plugin/addons/godot_mcp addons/godot_mcp
```
