# Godot MCP

**The standard [Model Context Protocol](https://modelcontextprotocol.io/) implementation for the [Godot Engine](https://godotengine.org/).**

Godot MCP connects any AI assistant — Claude, Cursor, ChatGPT, Gemini, Copilot, and others — to the Godot editor and runtime through a single, community-maintained MCP server.

> **Status:** v1.0.0 — 64 tools, 12 families, generic object model, event streaming, extension SDK.

---

## Table of Contents

- [Why Godot MCP?](#why-godot-mcp)
- [Features](#features)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Design Principles](#design-principles)
- [Project Structure](#project-structure)
- [Tool Reference](#tool-reference)
- [Events](#events)
- [Permissions](#permissions)
- [Compatibility](#compatibility)
- [Tech Stack](#tech-stack)
- [Roadmap](#roadmap)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

For architecture and protocol details, see **[docs/DESIGN.md](docs/DESIGN.md)**.

---

## Why Godot MCP?

Godot has no widely adopted, community-standard MCP server today — unlike ecosystems such as VS Code or Unity. Most integrations are one-off plugins tied to a specific AI client or IDE.

Godot MCP fills that gap by providing:

- **One protocol** — any MCP-compatible AI client can connect.
- **One install path** — no per-IDE configuration.
- **One tool surface** — stable, versioned tools instead of raw Godot APIs.

Target experience:

```text
brew install godot-mcp
        ↓
Open Godot → Enable Plugin
        ↓
AI can use Godot
```

---

## Features

| Capability | Description |
|------------|-------------|
| **AI-agnostic** | Works with any MCP client; not tied to Cursor, Claude, or a specific IDE |
| **Tool-based API** | High-level tools (`create_node`, `run_project`, …) — AI clients do not need Godot API knowledge |
| **Thin Godot plugin** | Plugin bridges Godot ↔ server; business logic lives in the MCP server |
| **Event streaming** | Server pushes editor/runtime events; no polling required |
| **Permission gates** | Destructive actions require user confirmation |
| **Extensible** | Third parties can add tools via SDK without modifying core |
| **Versioned tools** | Stable contracts so AI clients do not break on upgrades |

---

## Quick Start

> Available after v0.1 is released.

See [docs/HOMEBREW.md](docs/HOMEBREW.md) for Homebrew installation details.

```bash
# Install the CLI (Homebrew)
brew tap triforge0/godot-mcp https://github.com/triforge0/godot-mcp
brew install godot-mcp

# Or build from source
git clone https://github.com/triforge0/godot-mcp.git
cd godot-mcp && make build

# Start the MCP server
godot-mcp start

# Verify the setup
godot-mcp doctor
```

Then in Godot:

1. Open **Project → Project Settings → Plugins**
2. Enable **Godot MCP**
3. Configure your AI client to use the `godot-mcp` MCP server

See [docs.godot-mcp.org](https://docs.godot-mcp.org) for client-specific setup (Cursor, Claude Desktop, etc.).

---

## Architecture

```text
┌─────────────────────────────────────────┐
│  AI Clients                             │
│  Claude · Cursor · ChatGPT · Gemini …   │
└──────────────────┬──────────────────────┘
                   │ MCP (stdio / HTTP)
                   ▼
┌─────────────────────────────────────────┐
│  Godot MCP Server          (Go binary)  │
│  ┌───────────────────────────────────┐  │
│  │ Tool Registry · Session · Events  │  │
│  │ Dispatcher · Permissions · Cache  │  │
│  │ Logging · Versioning              │  │
│  └───────────────────────────────────┘  │
└──────────────────┬──────────────────────┘
                   │ WebSocket / JSON-RPC
                   ▼
┌─────────────────────────────────────────┐
│  Godot Plugin              (GDScript)   │
│  Scene · Node · Resource · Runtime      │
│  Editor · Filesystem · UndoRedo         │
└─────────────────────────────────────────┘
```

### Modules

| Module | Package | Role |
|--------|---------|------|
| **Core** | `godot-mcp-core` | MCP protocol, registry, dispatcher, session — no Godot dependency |
| **Server** | `godot-mcp-server` | Standalone process; implements all tools |
| **Plugin** | `godot-mcp-plugin` | Godot editor addon; thin bridge to the server |
| **CLI** | `godot-mcp` | `start`, `stop`, `doctor`, `version` |
| **SDK** | `sdk-go`, `sdk-python`, … | Build custom tools and extensions |

---

## Design Principles

### 1. Independent MCP server

The server speaks only the MCP protocol. It does not depend on Cursor, Claude, ChatGPT, or any specific IDE.

### 2. Thin Godot plugin

```text
Godot API  →  JSON-RPC  →  MCP Server
```

The plugin contains no business logic. This keeps it easy to maintain, reduces bugs, and improves compatibility across Godot versions.

### 3. Tools, not APIs

Expose named tools with clear inputs and outputs. AI clients call `save_scene()` — not `EditorInterface.save_scene()`.

### 4. Stable · Extensible · Safe

1. **Stable** — versioned tool contracts
2. **Extensible** — third-party tools via SDK
3. **Safe** — permission prompts for destructive operations

---

## Project Structure

```text
godot-mcp/
├── cmd/godot-mcp/
├── internal/
│   ├── mcp/           # server, registry, protocol, dispatcher
│   ├── bridge/        # jsonrpc.go, websocket.go
│   ├── tools/         # 12 families (project, scene, node, …)
│   ├── permission/
│   ├── events/
│   ├── godot/         # generic Client.Call(RPC)
│   ├── domain/        # GodotObject model
│   └── config/
├── plugin/addons/godot_mcp/
│   ├── bridge.gd
│   ├── rpc_router.gd
│   └── adapters/      # 12 thin Godot API adapters
├── spec/tools/
├── sdk/go/godotmcp/
└── docs/
```

Long term, core MCP infrastructure may be extracted into a reusable `mcp-core` package for other creative tools (Blender, Tiled, etc.).

---

## Tool Reference

Tools are grouped by domain. Full schemas will be published at [docs.godot-mcp.org](https://docs.godot-mcp.org).

### Editor

`open_project` · `close_project` · `save_project` · `reload_project`

### Scene

`open_scene` · `save_scene` · `scene_tree` · `duplicate_scene`

### Node

`create_node` · `delete_node` · `rename_node` · `move_node` · `set_property` · `get_property`

### Resource

`create_material` · `import_texture` · `load_scene` · `load_resource`

### Runtime

`run_project` · `stop_project` · `pause_project` · `resume_project` · `take_screenshot` · `simulate_input`

### Debug

`get_errors` · `get_console` · `get_profiler` · `get_fps`

### Script

`read_script` · `edit_script` · `attach_script` · `format_script` · `run_gdscript`

---

## Events

The server pushes events to connected AI clients — no polling required.

| Event | Trigger |
|-------|---------|
| `scene_changed` | Active scene modified |
| `node_selected` | Selection changed in editor |
| `project_opened` | Project loaded |
| `project_saved` | Project saved |
| `game_started` | Play mode started |
| `game_stopped` | Play mode stopped |

---

## Permissions

Destructive operations and `script_execute` show a native Godot editor dialog before running:

```text
AI requests: node_delete
        ↓
┌──────────────────────────────┐
│  Godot MCP — Permission      │
│  Allow Once | Allow Always   │
│  Deny                        │
└──────────────────────────────┘
```

"Allow Always" is cached for the current editor session. Set `GODOT_MCP_ALLOW_DESTRUCTIVE=1` or `GODOT_MCP_ALLOW_SCRIPT_EXEC=1` to skip dialogs during development.

---

## Compatibility

| Godot Version | Support |
|---------------|---------|
| 4.3 – 4.7 | ✅ Planned |
| 3.x | ❌ Not supported |

---

## Tech Stack

### MCP Server — Go

Go is the recommended server language:

- **Single binary** — no Java, Node, or Python runtime required
- **Cross-platform** — Windows, macOS (Intel & Apple Silicon), Linux
- **Strong fit** — JSON-RPC, WebSocket, IPC, filesystem, and process management
- **Simple contributor setup** — `git clone && go run .`

Suggested libraries:

| Purpose | Library |
|---------|---------|
| CLI | [`cobra`](https://github.com/spf13/cobra) |
| Logging | [`slog`](https://pkg.go.dev/log/slog) or [`zap`](https://github.com/uber-go/zap) |
| WebSocket | [`gorilla/websocket`](https://github.com/gorilla/websocket) |
| JSON | `encoding/json` |
| MCP | Official Go SDK when available; otherwise implement per MCP spec |

### Godot Plugin — GDScript

- **GDScript** (default) — best community fit, easy to install
- **C#** — only if special performance or integration needs arise

All business logic stays in the Go server; the plugin is a bridge only.

---

## Roadmap

| Version | Focus | Status |
|---------|-------|--------|
| **v0.1** | Connectivity | ✅ `ping`, `scene_tree`, `open_scene`, `save_scene`, `run_project` |
| **v0.2** | Editing | ✅ Node CRUD, `set_property`, `get_property` |
| **v0.3** | Runtime | ✅ Screenshot, logs, FPS, errors, inspector |
| **v0.4** | AI Features | ✅ Events, playtesting, input simulation, asset search |
| **v1.0** | Ecosystem | ✅ SDK, semver, marketplace, CI (Win/macOS/Linux) |

---

## Documentation

| Resource | Description |
|----------|-------------|
| [docs/TOOLS.md](docs/TOOLS.md) | Complete v1.0 tool reference |
| [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) | Build, install plugin, configure AI client |
| [docs/HOMEBREW.md](docs/HOMEBREW.md) | Homebrew tap and plugin install path |
| [docs/TESTING.md](docs/TESTING.md) | Unit, mock e2e, and real Godot integration tests |
| [docs/DESIGN.md](docs/DESIGN.md) | Architecture, protocol, permissions, and design decisions |
| [CONTRIBUTING.md](CONTRIBUTING.md) | How to contribute code, tools, and docs |
| [docs.godot-mcp.org](https://docs.godot-mcp.org) | Public docs site (Getting Started, Tool Reference, SDK, Examples) |

---

## Contributing

Contributions are welcome — see **[CONTRIBUTING.md](CONTRIBUTING.md)** for setup, PR guidelines, and the RFC process for significant design changes.

Example project: [examples/demo/](examples/demo/) — open in Godot 4.3+ to test integrations quickly.

---

## License

[MIT License](LICENSE) — Copyright (c) 2026 Godot MCP Contributors
