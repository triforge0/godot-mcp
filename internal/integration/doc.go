//go:build !integration

// Package integration contains end-to-end tests against a real Godot editor.
// Run them with:
//
//	go test -tags=integration -timeout 10m ./internal/integration/...
//
// Requires Godot 4.3+ (GODOT_BIN or godot on PATH).
package integration
