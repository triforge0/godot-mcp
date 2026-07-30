# Contributing to Godot MCP

Thank you for your interest in contributing to Godot MCP. This project aims to become the community-standard MCP implementation for the Godot Engine, and every contribution helps.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Ways to Contribute](#ways-to-contribute)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Design Changes and RFCs](#design-changes-and-rfcs)
- [Adding Tools](#adding-tools)
- [Reporting Issues](#reporting-issues)
- [Questions](#questions)

---

## Code of Conduct

Be respectful, constructive, and inclusive. Harassment, discrimination, and bad-faith behavior are not tolerated. Maintainers may remove contributions or restrict participation when necessary to keep the community healthy.

---

## Ways to Contribute

You do not need to write code to help:

- **Report bugs** — include reproduction steps and environment details
- **Suggest features** — explain the use case and which tools or events are affected
- **Improve documentation** — README, `docs/`, examples, and tool schemas
- **Review pull requests** — design feedback and test coverage are valuable
- **Write code** — server, plugin, CLI, SDK, tests, and examples

---

## Getting Started

> The project is in early design phase. Setup instructions will be updated as implementation lands.

### Prerequisites

| Component | Requirement |
|-----------|-------------|
| MCP Server | Go 1.22+ |
| Godot Plugin | Godot 4.3+ |
| AI client | Any MCP-compatible client (Cursor, Claude Desktop, etc.) |

### Clone and build (once code is available)

```bash
git clone https://github.com/your-org/godot-mcp.git
cd godot-mcp

# Build the CLI
go build -o bin/godot-mcp ./cmd/godot-mcp

# Run diagnostics
./bin/godot-mcp doctor
```

### Enable the Godot plugin

1. Copy or symlink `plugin/addons/godot_mcp/` into your project's `addons/` directory
2. Open **Project → Project Settings → Plugins**
3. Enable **Godot MCP**

---

## Development Workflow

1. **Check existing issues** — avoid duplicate work
2. **Open an issue or RFC** — for non-trivial changes, discuss before coding
3. **Fork and branch** — use descriptive branch names (`feat/scene-tree-tool`, `fix/websocket-reconnect`)
4. **Make focused changes** — one logical change per pull request
5. **Test locally** — verify server, plugin, and at least one MCP client
6. **Open a pull request** — fill in the template and link related issues

### Branch naming

```text
feat/<short-description>
fix/<short-description>
docs/<short-description>
refactor/<short-description>
```

---

## Pull Request Guidelines

Before submitting:

- [ ] Changes match the [design principles](docs/DESIGN.md)
- [ ] New tools include schema definitions in `spec/`
- [ ] Destructive tools integrate with the permission system
- [ ] Documentation is updated (README, `docs/`, or tool reference)
- [ ] Commits are clear and logically grouped

### PR description should include

1. **What** changed
2. **Why** it was needed
3. **How** to test it
4. **Screenshots or logs** when relevant (UI prompts, CLI output, MCP responses)

Maintainers may request changes. Once approved, your PR will be squashed or merged according to repository policy.

---

## Design Changes and RFCs

Significant changes require discussion before implementation:

- New tool categories or breaking tool schema changes
- Transport or protocol changes between server and plugin
- Permission model changes
- Extraction of `mcp-core` as a standalone framework

### RFC process

1. Open an issue with the `RFC` label, or add a document under `spec/rfc/`
2. Describe the problem, proposed solution, alternatives, and migration path
3. Allow at least one week for community feedback before starting large implementations
4. Link the accepted RFC in your pull request

See [docs/DESIGN.md](docs/DESIGN.md) for the current architecture and design rationale.

---

## Adding Tools

Tools are the public contract between AI clients and Godot. Follow these rules:

1. **High-level, not raw API** — expose intent (`save_scene`), not Godot internals
2. **Stable schemas** — define inputs/outputs in `spec/tools/<category>/<tool>.json`
3. **Idempotent when possible** — safe to retry where it makes sense
4. **Permission-aware** — mark destructive or sensitive tools in the permission registry
5. **Plugin stays thin** — implement orchestration in the Go server; plugin executes Godot calls

### Tool layout

```text
tools/
└── scene/
    └── save_scene.go      # Server-side handler

spec/
└── tools/
    └── scene/
        └── save_scene.json # JSON Schema for MCP tool definition

plugin/addons/godot_mcp/
└── handlers/
    └── scene.gd            # Godot-side execution only
```

---

## Reporting Issues

Include as much detail as possible:

```markdown
**Description**
What happened vs. what you expected.

**Environment**
- OS:
- Godot version:
- godot-mcp version:
- AI client (Cursor, Claude Desktop, etc.):

**Steps to reproduce**
1.
2.
3.

**Logs**
Paste relevant server, plugin, or MCP client logs.
```

Security vulnerabilities should **not** be reported in public issues. Email maintainers directly (contact TBD).

---

## Questions

- **Architecture and design** — read [docs/DESIGN.md](docs/DESIGN.md)
- **Tool reference** — [docs.godot-mcp.org](https://docs.godot-mcp.org) (when published)
- **General discussion** — GitHub Discussions (when enabled)

Thank you for helping build the standard MCP layer for Godot.
