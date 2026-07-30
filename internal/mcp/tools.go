package mcp

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterTools(mcpServer *mcpsdk.Server, registry *Registry, dispatcher *Dispatcher) {
	dispatch := func(ctx context.Context, def ToolDefinition, params any) (*mcpsdk.CallToolResult, any, error) {
		result, err := dispatcher.AuthorizeAndCall(ctx, def, params)
		if err != nil {
			return &mcpsdk.CallToolResult{
				Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: errorString(err)}},
				IsError: true,
			}, nil, nil
		}
		text, err := FormatResult(result)
		if err != nil {
			return &mcpsdk.CallToolResult{
				Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: errorString(err)}},
				IsError: true,
			}, nil, nil
		}
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: text}},
		}, result, nil
	}

	for _, def := range registry.All() {
		def.MCPRegister(mcpServer, def, dispatch)
	}
}
