package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/godot-mcp/godot-mcp/internal/events"
	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/permission"
)

type Dispatcher struct {
	client *godot.Client
	gate   *permission.Gate
	events *events.Bus
}

func NewDispatcher(client *godot.Client, gate *permission.Gate, bus *events.Bus) *Dispatcher {
	return &Dispatcher{client: client, gate: gate, events: bus}
}

func (d *Dispatcher) AuthorizeAndCall(ctx context.Context, def ToolDefinition, params any) (any, error) {
	if err := d.gate.Authorize(ctx, def.Permission, def.Name, def.RPCMethod); err != nil {
		return nil, err
	}
	result, err := def.Handler(ctx, d.client, params)
	if err != nil {
		return nil, err
	}
	d.events.Publish("tool.executed", map[string]string{
		"tool": def.Name,
		"rpc":  def.RPCMethod,
	})
	return result, nil
}

func FormatResult(value any) (string, error) {
	if value == nil {
		return "null", nil
	}
	formatted, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

func errorString(err error) string {
	return fmt.Sprintf("Error: %s", err.Error())
}
