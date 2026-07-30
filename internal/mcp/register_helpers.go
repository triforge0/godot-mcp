package mcp

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type DispatchFunc func(ctx context.Context, def ToolDefinition, params any) (*mcpsdk.CallToolResult, any, error)

func RegisterNoArgs(server *mcpsdk.Server, def ToolDefinition, dispatch DispatchFunc) {
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        def.Name,
		Description: def.Description,
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, _ struct{}) (*mcpsdk.CallToolResult, any, error) {
		return dispatch(ctx, def, nil)
	})
}

func RegisterToolWithArgs[T any](server *mcpsdk.Server, def ToolDefinition, dispatch DispatchFunc, mapParams func(T) any) {
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        def.Name,
		Description: def.Description,
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, args T) (*mcpsdk.CallToolResult, any, error) {
		return dispatch(ctx, def, mapParams(args))
	})
}
