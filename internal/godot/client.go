package godot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/godot-mcp/godot-mcp/internal/bridge"
)

// Client is the generic Godot adapter — all tools call RPC methods through it.
type Client struct {
	caller bridge.Caller
}

func NewClient(caller bridge.Caller) *Client {
	return &Client{caller: caller}
}

func (c *Client) Connected() bool {
	return c.caller.Connected()
}

// Call invokes any Godot RPC method and decodes into dest when provided.
func (c *Client) Call(ctx context.Context, method string, params, dest any) (any, error) {
	raw, err := c.caller.Call(ctx, method, params)
	if err != nil {
		return nil, err
	}
	if dest != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, dest); err != nil {
			return nil, fmt.Errorf("decode %s: %w", method, err)
		}
		return dest, nil
	}
	if len(raw) == 0 {
		return nil, nil
	}
	var generic any
	if err := json.Unmarshal(raw, &generic); err != nil {
		return string(raw), nil
	}
	return generic, nil
}

// CallMap is a convenience for untyped RPC results.
func (c *Client) CallMap(ctx context.Context, method string, params any) (map[string]any, error) {
	var result map[string]any
	_, err := c.Call(ctx, method, params, &result)
	return result, err
}
