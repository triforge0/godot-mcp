# Getting Started

This guide walks through running **godot-mcp v0.1** locally.

## Prerequisites

- Go 1.25+ (for building the MCP server)
- Godot 4.3+
- An MCP-compatible AI client (Cursor, Claude Desktop, etc.)

## 1. Build the MCP server

```bash
git clone https://github.com/your-org/godot-mcp.git
cd godot-mcp
make build
```

The binary is written to `bin/godot-mcp`.

## 2. Install the Godot plugin

**Option A — use the demo project**

```bash
# Open examples/demo/ in Godot 4.3+
```

The demo includes a symlink to the plugin. See [examples/demo/README.md](../examples/demo/README.md).

**Option B — add to your own project**

Copy or symlink the addon into your Godot project:

```bash
cp -R plugin/addons/godot_mcp /path/to/your/project/addons/
```

Or clone this repo and symlink:

```bash
ln -s /path/to/godot-mcp/plugin/addons/godot_mcp /path/to/your/project/addons/godot_mcp
```

Then in Godot:

1. **Project → Project Settings → Plugins**
2. Enable **Godot MCP**

The plugin connects to `ws://127.0.0.1:6505/ws` by default. Override with the `GODOT_MCP_BRIDGE_URL` environment variable if needed.

## 3. Configure your AI client

### Cursor

Add to `.cursor/mcp.json` (project) or `~/.cursor/mcp.json` (global):

```json
{
  "mcpServers": {
    "godot-mcp": {
      "command": "/absolute/path/to/godot-mcp/bin/godot-mcp",
      "args": ["start"]
    }
  }
}
```

Restart Cursor after saving. The MCP server starts automatically when Cursor connects.

### Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "godot-mcp": {
      "command": "/absolute/path/to/godot-mcp/bin/godot-mcp",
      "args": ["start"]
    }
  }
}
```

## 4. Verify connectivity

With the MCP server running (via your AI client) and Godot open with the plugin enabled:

```bash
make doctor
```

Expected output when everything is connected:

```text
godot-mcp 1.0.0
bridge address: 127.0.0.1:6505
bridge: reachable (connected)
plugin: connected
```

## 5. Try the v0.1 tools

| Tool | Description |
|------|-------------|
| `ping` | Check plugin connectivity and Godot version |
| `scene_tree` | Get the current scene tree |
| `open_scene` | Open a scene by `res://` path |
| `save_scene` | Save the active scene |
| `run_project` | Start play mode |

Ask your AI assistant, for example:

> Use the `ping` tool to check if Godot is connected.

> Call `scene_tree` and summarize the node hierarchy.

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| `godot plugin is not connected` | Open Godot, enable the plugin, ensure port 6505 is free |
| `bridge: not reachable` | Confirm your AI client has started `godot-mcp start` |
| Plugin won't connect | Check firewall; set `GODOT_MCP_BRIDGE_ADDR` if using a custom port |

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GODOT_MCP_BRIDGE_ADDR` | `127.0.0.1:6505` | WebSocket listen address (server side) |
| `GODOT_MCP_BRIDGE_URL` | `ws://127.0.0.1:6505/ws` | WebSocket URL (Godot plugin side) |
| `GODOT_MCP_ALLOW_DESTRUCTIVE` | (unset) | Skip destructive permission dialogs when `1` |
| `GODOT_MCP_ALLOW_SCRIPT_EXEC` | (unset) | Skip `script_execute` permission dialog when `1` |
