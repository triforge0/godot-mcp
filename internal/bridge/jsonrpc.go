package bridge

import (
	"context"
	"encoding/json"
)

const JSONRPCVersion = "2.0"

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type EventParams struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

const EventMethod = "event"

// Caller is implemented by WebSocket and test mocks.
type Caller interface {
	Call(ctx context.Context, method string, params any) (json.RawMessage, error)
	Connected() bool
}
