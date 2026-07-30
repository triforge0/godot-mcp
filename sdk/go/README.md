# Godot MCP Go SDK

Register custom MCP tools in third-party Go modules:

```go
package mytools

import (
    "context"

    "github.com/godot-mcp/godot-mcp/sdk/go/godotmcp"
)

func Register(r *godotmcp.Registry) {
    r.Register(godotmcp.ToolDefinition{
        Name:        "my_custom_tool",
        Description: "Does something useful in Godot.",
        RPCMethod:   "custom.my_tool",
        Permission:  godotmcp.LevelRead,
        Handler: func(ctx context.Context, client *godotmcp.Client, _ any) (any, error) {
            return client.Ping(ctx)
        },
    })
}
```

Wire into your server by calling `mytools.Register(registry)` before starting MCP.

See [CONTRIBUTING.md](../../CONTRIBUTING.md) and [marketplace/README.md](../../marketplace/README.md).
