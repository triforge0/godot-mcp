# Godot MCP v1.0 — Tool Matrix

**64 tools** across **12 families**. MCP tool names use `snake_case`; RPC methods use `domain.action`.

## Architecture

```text
MCP Tool  →  Permission  →  Handler  →  godot.Client.Call(RPC)  →  Plugin Adapter
```

Plugin = thin Godot API adapter (12 families). No business logic in GDScript.

## Families

### Project (3)
| MCP Tool | RPC |
|----------|-----|
| `project_info` | `project.info` |
| `project_settings` | `project.settings` |
| `project_reload` | `project.reload` |

### Scene (7)
| MCP Tool | RPC |
|----------|-----|
| `scene_list` | `scene.list` |
| `scene_current` | `scene.current` |
| `scene_open` | `scene.open` |
| `scene_save` | `scene.save` |
| `scene_create` | `scene.create` |
| `scene_reload` | `scene.reload` |
| `scene_close` | `scene.close` |

### Node (10)
| MCP Tool | RPC |
|----------|-----|
| `node_list` | `node.list` |
| `node_create` | `node.create` |
| `node_delete` | `node.delete` |
| `node_rename` | `node.rename` |
| `node_move` | `node.move` |
| `node_duplicate` | `node.duplicate` |
| `node_get_property` | `node.get_property` |
| `node_set_property` | `node.set_property` |
| `node_children` | `node.children` |
| `node_parent` | `node.parent` |

### Resource (6)
| MCP Tool | RPC |
|----------|-----|
| `resource_load` | `resource.load` |
| `resource_save` | `resource.save` |
| `resource_create` | `resource.create` |
| `resource_inspect` | `resource.inspect` |
| `resource_list` | `resource.list` |
| `resource_delete` | `resource.delete` |

### FileSystem (6)
| MCP Tool | RPC |
|----------|-----|
| `file_list` | `filesystem.list` |
| `file_read` | `filesystem.read` |
| `file_write` | `filesystem.write` |
| `file_create` | `filesystem.create` |
| `file_delete` | `filesystem.delete` |
| `folder_create` | `filesystem.mkdir` |

### Script (6)
| MCP Tool | RPC |
|----------|-----|
| `script_read` | `script.read` |
| `script_create` | `script.create` |
| `script_update` | `script.update` |
| `script_attach` | `script.attach` |
| `script_detach` | `script.detach` |
| `script_execute` | `script.execute` |

### Runtime (6)
| MCP Tool | RPC |
|----------|-----|
| `runtime_run` | `runtime.run` |
| `runtime_stop` | `runtime.stop` |
| `runtime_pause` | `runtime.pause` |
| `runtime_resume` | `runtime.resume` |
| `runtime_status` | `runtime.status` |
| `runtime_restart` | `runtime.restart` |

### Screenshot (2)
| MCP Tool | RPC |
|----------|-----|
| `screenshot_capture` | `screenshot.capture` |
| `screenshot_viewport` | `screenshot.viewport` |

### Console / Debug (5)
| MCP Tool | RPC |
|----------|-----|
| `console_logs` | `console.logs` |
| `console_clear` | `console.clear` |
| `errors_list` | `errors.list` |
| `errors_clear` | `errors.clear` |
| `profiler_stats` | `profiler.stats` |

### Editor (5)
| MCP Tool | RPC |
|----------|-----|
| `editor_selection` | `editor.selection` |
| `editor_focus` | `editor.focus` |
| `editor_undo` | `editor.undo` |
| `editor_redo` | `editor.redo` |
| `editor_refresh` | `editor.refresh` |

### Reflection (4)
| MCP Tool | RPC |
|----------|-----|
| `object_inspect` | `object.inspect` |
| `class_inspect` | `class.inspect` |
| `property_list` | `property.list` |
| `method_list` | `method.list` |

### Skills (4)
| MCP Tool | RPC |
|----------|-----|
| `skill_create_player` | `skill.create_player` |
| `skill_create_scene` | `skill.create_scene` |
| `skill_optimize_project` | `skill.optimize_project` |
| `skill_analyze_error` | `skill.analyze_error` |

## Generic Object Model

All node/resource tools use `domain.GodotObject`:

```go
type GodotObject struct {
    Class      string
    Path       string
    Name       string
    Properties map[string]any
    Children   []GodotObject
}
```

Property access is reflection-based — no hard-coded node types.

## Events

| Event | Source |
|-------|--------|
| `node.selected` | Editor selection |
| `scene.changed` | Scene changes |
| `runtime.started` | Play mode |
| `runtime.stopped` | Stop play |
| `runtime.error` | Error buffer |
| `tool.executed` | Go server |

Flow: **Godot → WebSocket event → EventBus → MCP notification → AI**

## Permissions

| Level | Examples |
|-------|----------|
| read | list, inspect, logs |
| write | create, set_property |
| execute | runtime.run, script.execute |
| destructive | delete_node, file_delete |

Dev overrides:
- `GODOT_MCP_ALLOW_DESTRUCTIVE=1`
- `GODOT_MCP_ALLOW_SCRIPT_EXEC=1`
