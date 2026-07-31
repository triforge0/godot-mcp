# Setup: Connecting godot-mcp on macOS and Windows

This guide walks through getting `godot-mcp` running end-to-end: CLI installed, Godot plugin enabled, and an MCP client (Claude Code, Claude Desktop, or Cursor) talking to a live Godot editor.

Recap of the moving pieces:

```
AI client  <--MCP over stdio-->  godot-mcp (CLI)  <--WebSocket-->  Godot editor (plugin enabled)
```

The AI client spawns `godot-mcp start`, which opens a WebSocket bridge on `127.0.0.1:6505` by default. The Godot plugin connects out to that bridge when enabled in the editor. Nothing works until **both** the CLI is running (via your AI client) **and** the plugin is enabled in an open Godot project.

---

## macOS

### 1. Install the CLI

**Homebrew (recommended)** — builds from the tagged release source and also bundles the Godot plugin + demo project:

```bash
brew tap triforge0/godot-mcp https://github.com/triforge0/godot-mcp
brew install godot-mcp
```

If Homebrew refuses the tap as untrusted, trust it explicitly (Homebrew's tap-trust security prompt):

```bash
brew trust triforge0/godot-mcp
```

Verify:

```bash
godot-mcp version
```

**Or: download the prebuilt binary** from [GitHub Releases](https://github.com/triforge0/godot-mcp/releases) (`godot-mcp-darwin-arm64` for Apple Silicon, `godot-mcp-darwin-amd64` for Intel). This gets you only the binary — you still need the plugin separately (see step 2, option B). macOS Gatekeeper quarantines files downloaded this way, so clear that before running it:

```bash
chmod +x godot-mcp-darwin-arm64
xattr -d com.apple.quarantine godot-mcp-darwin-arm64
mv godot-mcp-darwin-arm64 /usr/local/bin/godot-mcp   # or anywhere on your PATH
```

### 2. Install the Godot plugin

**Option A — Homebrew install**: the formula bundles the plugin under Homebrew's share directory. Symlink it into your project:

```bash
ln -s "$(brew --prefix godot-mcp)/share/godot-mcp/plugin/addons/godot_mcp" \
  /path/to/your/project/addons/godot_mcp
```

(If `addons/` doesn't exist yet in your project, create it first: `mkdir -p /path/to/your/project/addons`.)

**Option B — binary-only install**: clone the repo (or download the source tarball from the release page) to get `plugin/addons/godot_mcp/`, then symlink or copy it the same way:

```bash
git clone https://github.com/triforge0/godot-mcp.git
ln -s "$(pwd)/godot-mcp/plugin/addons/godot_mcp" /path/to/your/project/addons/godot_mcp
```

### 3. Enable the plugin in Godot

Open your project in Godot 4.3+, then **Project → Project Settings → Plugins** → enable **Godot MCP**.

### 4. Point your AI client at the CLI

**Claude Code:**

```bash
claude mcp add godot-mcp -- godot-mcp start
# or, to make it available in every project:
claude mcp add godot-mcp --scope user -- godot-mcp start
```

Restart Claude Code (MCP servers are loaded at session start, so a server added mid-session won't show up until you reopen).

**Claude Desktop** — edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "godot-mcp": {
      "command": "/opt/homebrew/bin/godot-mcp",
      "args": ["start"]
    }
  }
}
```

(Use the actual path from `which godot-mcp` — Homebrew on Apple Silicon is `/opt/homebrew/bin/godot-mcp`; on Intel Macs it's usually `/usr/local/bin/godot-mcp`.)

**Cursor** — same shape, in `.cursor/mcp.json` (project) or `~/.cursor/mcp.json` (global):

```json
{
  "mcpServers": {
    "godot-mcp": {
      "command": "/opt/homebrew/bin/godot-mcp",
      "args": ["start"]
    }
  }
}
```

### 5. Verify

With Godot open and the plugin enabled, run:

```bash
godot-mcp doctor
```

Expect:

```
godot-mcp 1.0.0
tools registered: 64
bridge address: 127.0.0.1:6505

bridge: reachable (connected)
Status: READY
plugin version: 1.0.0
godot version: 4.x.x
project: <your project name>
```

---

## Windows

### 1. Install the CLI

There's no Homebrew on Windows, so install one of these ways:

**Download the prebuilt binary** from [GitHub Releases](https://github.com/triforge0/godot-mcp/releases): `godot-mcp-windows-amd64.exe`. Rename it to `godot-mcp.exe` and put it somewhere on your `PATH` (or just reference the full path in your AI client config — see step 4).

Windows SmartScreen will likely flag the unsigned `.exe` on first run ("Windows protected your PC"). Either click **More info → Run anyway**, or unblock it ahead of time in PowerShell:

```powershell
Unblock-File -Path .\godot-mcp.exe
```

**Or build from source** (needs [Go 1.25+](https://go.dev/dl/)):

```powershell
git clone https://github.com/triforge0/godot-mcp.git
cd godot-mcp
go build -o godot-mcp.exe ./cmd/godot-mcp
```

Verify:

```powershell
.\godot-mcp.exe version
```

### 2. Get the Godot plugin

The release binary is CLI-only — the plugin isn't bundled with it on Windows. Get it from the source tree:

```powershell
git clone https://github.com/triforge0/godot-mcp.git
```

(No Go toolchain needed just for this — you only need the `plugin/` folder out of the clone.)

Copy it into your Godot project's `addons/` folder:

```powershell
mkdir -Force "C:\path\to\your\project\addons"
Copy-Item -Recurse ".\godot-mcp\plugin\addons\godot_mcp" "C:\path\to\your\project\addons\godot_mcp"
```

A real symlink (`mklink /D`) works too, but needs either an elevated (Admin) prompt or Developer Mode enabled (**Settings → Privacy & security → For developers → Developer Mode**). Copying is simpler and avoids that requirement — the plugin has no reason to stay in sync with the source tree, so a plain copy is fine.

### 3. Enable the plugin in Godot

Open your project in Godot 4.3+, then **Project → Project Settings → Plugins** → enable **Godot MCP**.

### 4. Point your AI client at the CLI

**Claude Code** (same command, from a shell that has `godot-mcp.exe` on `PATH`, or use the full path):

```powershell
claude mcp add godot-mcp -- godot-mcp.exe start
```

**Claude Desktop** — edit `%APPDATA%\Claude\claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "godot-mcp": {
      "command": "C:/path/to/godot-mcp.exe",
      "args": ["start"]
    }
  }
}
```

(Forward slashes work fine here even on Windows — no need to escape backslashes in the JSON.)

**Cursor** — same shape, in `.cursor/mcp.json` (project) or `%USERPROFILE%\.cursor\mcp.json` (global).

### 5. Verify

With Godot open and the plugin enabled:

```powershell
.\godot-mcp.exe doctor
```

Same expected output as the macOS section above.

---

## Troubleshooting (both platforms)

| Symptom | Fix |
|---|---|
| `Status: NOT READY` / `bridge: not reachable` | Your AI client hasn't started `godot-mcp start` yet — confirm the MCP server shows as connected in your client (e.g. `/mcp` in Claude Code), and that you restarted the client after adding it. |
| `Status: WAITING FOR PLUGIN` | Bridge is up but Godot hasn't connected. Open the project in Godot and confirm **Godot MCP** is enabled under Project Settings → Plugins. |
| Plugin never connects even though enabled | Check nothing else is bound to port `6505` (`lsof -i :6505` / `netstat -ano | findstr 6505`). Override the port on both sides with `GODOT_MCP_BRIDGE_ADDR` (server) and `GODOT_MCP_BRIDGE_URL` (plugin) if needed. |
| Tools missing in an already-open chat session | MCP servers are only loaded at session start. Restart your AI client/session after running `claude mcp add`. |
| Destructive tools keep prompting a dialog in Godot | Expected — that's the permission gate. Set `GODOT_MCP_ALLOW_DESTRUCTIVE=1` / `GODOT_MCP_ALLOW_SCRIPT_EXEC=1` in the environment `godot-mcp start` runs with to bypass during development. |

## Environment variables

| Variable | Default | Used by |
|---|---|---|
| `GODOT_MCP_BRIDGE_ADDR` | `127.0.0.1:6505` | CLI: address the bridge listens on |
| `GODOT_MCP_BRIDGE_URL` | `ws://127.0.0.1:6505/ws` | Godot plugin: URL it connects out to |
| `GODOT_MCP_ALLOW_DESTRUCTIVE` | unset | `1` skips the permission dialog for destructive tools |
| `GODOT_MCP_ALLOW_SCRIPT_EXEC` | unset | `1` skips the permission dialog for `script_execute` |
