# Godot MCP — Design Document

This document describes the architecture, protocol, and design decisions behind Godot MCP. It is intended for contributors and integrators who need more detail than the [README](../README.md).

---

## Table of Contents

- [Vision](#vision)
- [Goals and Non-Goals](#goals-and-non-goals)
- [System Overview](#system-overview)
- [Components](#components)
- [Communication Protocol](#communication-protocol)
- [Tool Model](#tool-model)
- [Event System](#event-system)
- [Permission Model](#permission-model)
- [Session and State](#session-and-state)
- [Versioning](#versioning)
- [Tech Stack Rationale](#tech-stack-rationale)
- [Future: mcp-core Extraction](#future-mcp-core-extraction)
- [Repository Layout](#repository-layout)
- [Roadmap](#roadmap)

---

## Vision

> **The standard MCP implementation for the Godot Engine.**

Godot lacks a community-recognized MCP server comparable to what exists for VS Code or Unity. Godot MCP provides a single, AI-agnostic integration layer so any MCP-compatible assistant can interact with the Godot editor and runtime through a stable set of tools.

Target user flow:

```text
brew install godot-mcp
        ↓
Open Godot → Enable Plugin
        ↓
AI can use Godot
```

No per-IDE configuration. No vendor lock-in.

---

## Goals and Non-Goals

### Goals

- **AI-agnostic** — works with Claude, Cursor, ChatGPT, Gemini, Copilot, and any MCP client
- **Simple install** — single binary + Godot addon
- **Stable tool contracts** — versioned schemas that AI clients can rely on
- **Extensible** — third parties can add tools via SDK
- **Safe by default** — permission gates for destructive operations
- **Event-driven** — server pushes editor/runtime changes to AI clients

### Non-Goals

- **Godot 3 support** — Godot 4.3+ only
- **Replacing the Godot editor** — this is an integration layer, not a new IDE
- **Embedding AI inside Godot** — AI runs in external clients; Godot MCP is the bridge
- **Exposing raw Godot APIs over MCP** — tools are curated and high-level

---

## System Overview

```text
┌─────────────────────────────────────────┐
│  AI Clients                             │
│  Claude · Cursor · ChatGPT · Gemini …   │
└──────────────────┬──────────────────────┘
                   │ MCP (stdio / HTTP)
                   ▼
┌─────────────────────────────────────────┐
│  Godot MCP Server          (Go binary)  │
│                                         │
│  Tool Registry    Session Manager       │
│  Dispatcher       Event Bus             │
│  Permission Gate  Cache                 │
│  Logger           Version Manager       │
└──────────────────┬──────────────────────┘
                   │ WebSocket + JSON-RPC
                   ▼
┌─────────────────────────────────────────┐
│  Godot Plugin              (GDScript)   │
│                                         │
│  Scene · Node · Resource · Runtime      │
│  Editor · Filesystem · UndoRedo         │
└─────────────────────────────────────────┘
```

### Data flow for a tool call

```text
AI Client
  │  MCP: tools/call { name: "save_scene", arguments: { ... } }
  ▼
MCP Server
  │  validate schema → check permissions → dispatch
  ▼
WebSocket JSON-RPC
  │  { method: "scene.save", params: { ... } }
  ▼
Godot Plugin
  │  EditorInterface.save_scene() (or equivalent)
  ▼
Godot Editor
  │  result / error
  ▼
MCP Server → AI Client
```

The plugin never implements business logic. It translates JSON-RPC calls into Godot API invocations and returns structured results.

---

## Implementation Architecture

The v0.1 codebase follows a layered design optimized for community contributions:

```text
Tool (MCP name)
    ↓
Permission Gate
    ↓
Handler (internal/tools/<category>/)
    ↓
Godot Client (internal/godot/)
    ↓
Bridge RPC (internal/bridge/rpc/)
    ↓
WebSocket Transport (internal/bridge/websocket/)
    ↓
Godot Plugin
```

### Key packages

| Package | Role |
|---------|------|
| `internal/mcp/` | MCP server, tool registry, dispatcher |
| `internal/tools/` | One package per tool category; each tool registers via `Registry` |
| `internal/permissions/` | Permission gate (Tool → Permission → RPC) |
| `internal/events/` | Event bus for future MCP notifications |
| `internal/domain/` | Typed domain models (`SceneNode`, `PingResult`, …) |
| `internal/godot/` | Godot client with domain services (`Scene`, `Run`) |
| `internal/bridge/` | Transport layer — RPC types + WebSocket server |

Tool names and RPC methods are **decoupled** via `ToolDefinition`:

```go
registry.Register(mcp.ToolDefinition{
    Name:       "scene_tree",   // MCP tool name
    RPCMethod:  "scene.tree",   // Godot plugin RPC method
    Permission: permissions.LevelRead,
    Handler:    treeHandler,
})
```

Contributors add a tool by creating a package under `internal/tools/` and calling `Register()` in `internal/tools/register.go`.

---

## Components

### godot-mcp-core

Language-agnostic MCP infrastructure (may later become standalone `mcp-core`):

| Responsibility | Description |
|----------------|-------------|
| Protocol | MCP message handling (initialize, tools/list, tools/call, …) |
| Registry | Tool registration and lookup |
| Dispatcher | Route tool calls to handlers |
| Session | Per-client state and connection lifecycle |

No dependency on Godot.

### godot-mcp-server

The standalone process that AI clients connect to:

- Implements all Godot-specific tools
- Manages WebSocket connection to the Godot plugin
- Handles permissions, caching, logging, and events
- Ships as a single cross-platform binary

### godot-mcp-plugin

Godot editor addon (`plugin/addons/godot_mcp/`):

- Connects to the server via WebSocket on editor startup
- Executes Godot API calls requested by the server
- Emits editor/runtime events back to the server
- Written in **GDScript** by default; C# only if needed

### godot-mcp CLI

User-facing commands:

```bash
godot-mcp start     # Start the MCP server
godot-mcp stop      # Stop a running server
godot-mcp doctor    # Diagnose connectivity and configuration
godot-mcp version   # Print version info
```

### SDK

Language bindings for third-party tool authors:

| Language | Example API |
|----------|-------------|
| Go | `godotmcp.RegisterTool(...)` |
| Python | `@tool` decorator |
| Node.js | `registerTool(...)` |
| Java | `@GodotTool` annotation |

SDKs wrap the extension API so contributors can add tools without modifying core server code.

---

## Communication Protocol

### AI Client ↔ MCP Server

Standard [Model Context Protocol](https://modelcontextprotocol.io/) over:

- **stdio** — default for local AI clients (Cursor, Claude Desktop)
- **HTTP/SSE** — optional for remote or web-based clients

### MCP Server ↔ Godot Plugin

Custom JSON-RPC over WebSocket:

```json
// Request (server → plugin)
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "method": "node.create",
  "params": {
    "parent_path": "/root/Main",
    "type": "Sprite2D",
    "name": "Player"
  }
}

// Response (plugin → server)
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "result": {
    "path": "/root/Main/Player"
  }
}
```

Method naming convention: `<domain>.<action>` (e.g. `scene.open`, `runtime.screenshot`).

### Why WebSocket?

- Persistent connection for bidirectional communication
- Server can push events without polling
- Low latency for interactive editor operations
- Well-supported in Go and Godot

---

## Tool Model

### Principles

1. **Intent-based** — tools describe *what* to do, not *how* Godot implements it
2. **Schema-defined** — every tool has a JSON Schema in `spec/tools/`
3. **Composable** — small, focused tools preferred over monolithic ones
4. **Permission-tagged** — each tool declares its risk level

### Tool categories

| Category | Examples | Notes |
|----------|----------|-------|
| **Editor** | `open_project`, `save_project` | Project-level operations |
| **Scene** | `open_scene`, `scene_tree` | Scene management |
| **Node** | `create_node`, `set_property` | Scene tree manipulation |
| **Resource** | `import_texture`, `load_resource` | Asset operations |
| **Runtime** | `run_project`, `take_screenshot` | Play mode and simulation |
| **Debug** | `get_errors`, `get_console` | Diagnostics |
| **Script** | `read_script`, `edit_script` | GDScript operations |

Full tool list: see [README — Tool Reference](../README.md#tool-reference).

### Example tool schema

```json
{
  "name": "save_scene",
  "description": "Save the currently open scene to disk.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "Optional override path. Defaults to the scene's current path."
      }
    }
  },
  "annotations": {
    "permission": "write",
    "destructive": false,
    "idempotent": true
  }
}
```

---

## Event System

AI clients benefit from push notifications instead of polling the editor state.

### Event flow

```text
Godot Editor (selection change)
        ↓
Godot Plugin (emit event)
        ↓
MCP Server (Event Bus)
        ↓
AI Client (MCP notification)
```

### Planned events

| Event | Payload (summary) | Trigger |
|-------|-------------------|---------|
| `scene_changed` | `{ scene_path }` | Active scene modified |
| `node_selected` | `{ node_path }` | Editor selection changed |
| `project_opened` | `{ project_path }` | Project loaded |
| `project_saved` | `{ project_path }` | Project saved |
| `game_started` | `{}` | Entered play mode |
| `game_stopped` | `{}` | Exited play mode |
| `error_logged` | `{ message, source }` | Runtime or editor error |

Events are delivered to all connected MCP sessions subscribed to the relevant topic.

---

## Permission Model

Destructive or sensitive operations require explicit user approval in the Godot editor.

### Risk levels

| Level | Behavior | Examples |
|-------|----------|----------|
| **read** | Always allowed | `scene_tree`, `get_property`, `read_script` |
| **write** | Prompt on first use | `save_scene`, `set_property`, `edit_script` |
| **destructive** | Always prompt | `delete_node`, `delete_scene`, overwrite resources |
| **execute** | Always prompt | `run_gdscript`, `simulate_input` |

### Prompt flow

```text
AI requests: delete_node
        ↓
Server checks permission registry
        ↓
Plugin shows native Godot dialog:
  "AI wants to delete node 'Player'. Allow?"
  [ Once ] [ Always ] [ Never ] [ Cancel ]
        ↓
Decision cached per session / persisted per user preference
        ↓
Tool executes or returns permission_denied error
```

This mirrors the permission model used by VS Code extensions — familiar to developers and safe for AI-driven workflows.

---

## Session and State

| State | Owner | Lifetime |
|-------|-------|----------|
| MCP client session | Server | Duration of AI client connection |
| Godot plugin connection | Server + Plugin | Duration of editor session |
| Permission grants | Server (persisted) | Until revoked by user |
| Scene cache | Server | Invalidated on `scene_changed` events |
| Undo/redo stack | Godot | Managed by editor; plugin invokes via Godot API |

The server is stateless regarding scene content — Godot remains the source of truth.

---

## Versioning

### Server versioning

Semantic versioning: `MAJOR.MINOR.PATCH`

- **MAJOR** — breaking tool schema or protocol changes
- **MINOR** — new tools or backward-compatible features
- **PATCH** — bug fixes

### Tool schema versioning

Each tool schema includes a `version` field. Breaking changes require a new tool name or a major server bump with migration notes.

### Godot compatibility matrix

| Godot Version | godot-mcp Support |
|---------------|-------------------|
| 4.3 | ✅ Planned |
| 4.4 | ✅ Planned |
| 4.5 | ✅ Planned |
| 4.6 | ✅ Planned |
| 4.7 | ✅ Planned |
| 3.x | ❌ Not supported |

The plugin may use feature detection for minor API differences across Godot 4.x releases.

---

## Tech Stack Rationale

### Go for the MCP server

| Criterion | Why Go |
|-----------|--------|
| Distribution | Single static binary; no runtime install |
| Cross-platform | Easy builds for Windows, macOS, Linux |
| Concurrency | Native goroutines for WebSocket, events, screenshots |
| Ecosystem | Strong JSON-RPC, WebSocket, and CLI libraries |
| Contributor UX | `go run .` — no Maven, Gradle, or node_modules |

Suggested libraries:

| Purpose | Library |
|---------|---------|
| CLI | [cobra](https://github.com/spf13/cobra) |
| Logging | [slog](https://pkg.go.dev/log/slog) or [zap](https://github.com/uber-go/zap) |
| WebSocket | [gorilla/websocket](https://github.com/gorilla/websocket) |
| JSON | `encoding/json` |
| MCP | Official Go SDK when available |

### GDScript for the Godot plugin

- Default language for Godot community projects
- No extra toolchain for contributors
- Direct access to editor and runtime APIs
- C# reserved for exceptional performance needs

---

## Future: mcp-core Extraction

As the project matures, generic MCP infrastructure may be extracted:

```text
godot-mcp/
├── mcp-core/         # Reusable MCP framework (transport, session, registry)
├── godot-plugin/     # Godot-specific addon
└── godot-tools/      # Godot-specific tool implementations
```

`mcp-core` would handle transport, session management, registry, and dispatcher — reusable for other creative tools (Blender, Tiled, internal pipelines) without rewriting MCP plumbing.

Extraction is planned for post-v1.0 when the abstractions are proven stable.

---

## Repository Layout

```text
godot-mcp/
├── cmd/
│   └── godot-mcp/              # CLI entrypoint
├── internal/
│   ├── mcp/                    # MCP server, registry, dispatcher
│   ├── bridge/
│   │   ├── rpc/                # JSON-RPC types and Caller interface
│   │   └── websocket/          # WebSocket transport server
│   ├── godot/                  # Godot client (scene, runtime services)
│   ├── domain/                 # Domain models
│   ├── tools/                  # MCP tool packages (ping, scene, runtime, …)
│   ├── permissions/            # Permission gate
│   ├── events/                 # Event bus
│   ├── config/
│   └── version/
├── plugin/
│   └── addons/godot_mcp/       # Godot editor addon (thin bridge)
├── spec/
│   ├── tools/                  # MCP tool JSON Schemas
│   ├── rpc/                    # Godot plugin RPC method specs
│   └── rfc/                    # Design proposals
├── examples/
├── scripts/
├── docs/
└── website/
```

---

## Roadmap

| Version | Focus | Status |
|---------|-------|--------|
| **v0.1** | Connectivity | ✅ Shipped |
| **v0.2** | Editing | ✅ Shipped |
| **v0.3** | Runtime & Debug | ✅ Shipped |
| **v0.4** | AI Features | ✅ Shipped |
| **v1.0** | Ecosystem | ✅ SDK, semver, marketplace, CI |

### Differentiators at v1.0

1. **Stable** — versioned tool contracts with migration guides
2. **Extensible** — SDK for third-party tools without core patches
3. **Safe** — permission system for all sensitive operations

If executed well, Godot MCP becomes the de facto AI integration layer for the Godot ecosystem — one protocol, one tool set, every AI client.
