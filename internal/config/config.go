package config

import "os"

const DefaultBridgeAddr = "127.0.0.1:6505"

type Config struct {
	BridgeAddr string
}

func Load() Config {
	addr := os.Getenv("GODOT_MCP_BRIDGE_ADDR")
	if addr == "" {
		addr = DefaultBridgeAddr
	}
	return Config{BridgeAddr: addr}
}
