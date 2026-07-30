# Homebrew

Install the MCP server CLI with Homebrew:

```bash
brew tap triforge0/godot-mcp https://github.com/triforge0/godot-mcp
brew install godot-mcp
```

Verify:

```bash
godot-mcp version
godot-mcp doctor
```

## Godot plugin

The formula installs the editor plugin under Homebrew's share directory:

```bash
brew --prefix godot-mcp
# plugin at: $(brew --prefix)/share/godot-mcp/plugin/addons/godot_mcp
```

Symlink into your Godot project:

```bash
ln -s "$(brew --prefix godot-mcp)/share/godot-mcp/plugin/addons/godot_mcp" \
  /path/to/your/project/addons/godot_mcp
```

Then enable **Godot MCP** in **Project → Project Settings → Plugins**.

## Development (local formula)

From a clone of this repository:

```bash
brew install --formula ./Formula/godot-mcp.rb
```

## Releases

Pre-built binaries are attached to [GitHub Releases](https://github.com/triforge0/godot-mcp/releases). The Homebrew formula builds from source (HEAD/main or tagged releases when configured).
