package permission

import "context"

// Request is sent to the Godot plugin to show a permission dialog.
type Request struct {
	Tool    string         `json:"tool"`
	RPC     string         `json:"rpc"`
	Level   Level          `json:"level"`
	Details map[string]any `json:"details,omitempty"`
}

// Response is returned after the user approves or denies the action.
type Response struct {
	Approved bool   `json:"approved"`
	Remember string `json:"remember"` // "once", "always", or "none"
}

// Approver requests user approval from the Godot editor plugin.
type Approver interface {
	Connected() bool
	RequestPermission(ctx context.Context, req Request) (Response, error)
}
