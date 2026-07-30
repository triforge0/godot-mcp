# Godot MCP Marketplace

Community extensions for Godot MCP — additional tool packs for animation, tilemaps, shaders, multiplayer, and more.

## Publishing an extension

1. Implement tools using the [Go SDK](../sdk/go/README.md)
2. Add JSON schemas under `spec/tools/` in your repo
3. Open a PR to add an entry to [registry.json](registry.json)
4. Tag releases with [Semantic Versioning](https://semver.org/)

### Registry entry example

```json
{
  "name": "godot-mcp-tilemap",
  "version": "1.0.0",
  "author": "your-name",
  "description": "TileMap editing tools for Godot MCP",
  "repository": "https://github.com/you/godot-mcp-tilemap",
  "tools": ["create_tilemap_layer", "paint_tiles"],
  "categories": ["tilemap", "2d"]
}
```

## Categories

| Category | Examples |
|----------|----------|
| `animation` | Sprite animation, AnimationPlayer |
| `tilemap` | TileMapLayer, terrain painting |
| `shader` | Material and shader editing |
| `multiplayer` | Network, RPC debugging |
| `ui` | Control layout, theme tools |

## Official vs community

- **Core tools** ship with `godot-mcp` (connectivity, editing, runtime, debug)
- **Marketplace extensions** are optional add-ons maintained by the community
