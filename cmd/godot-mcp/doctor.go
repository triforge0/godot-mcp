package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/godot-mcp/godot-mcp/internal/config"
	"github.com/godot-mcp/godot-mcp/internal/version"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose Godot MCP connectivity",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := config.Load()

			fmt.Printf("godot-mcp %s\n", version.Version)
			fmt.Printf("bridge address: %s\n", cfg.BridgeAddr)

			client := &http.Client{Timeout: 2 * time.Second}
			resp, err := client.Get(fmt.Sprintf("http://%s/health", cfg.BridgeAddr))
			if err != nil {
				fmt.Println("bridge: not reachable (start `godot-mcp start` and enable the Godot plugin)")
				return nil
			}
			defer resp.Body.Close()

			var health map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
				return fmt.Errorf("decode health response: %w", err)
			}

			fmt.Printf("bridge: reachable (%s)\n", health["status"])
			if health["status"] != "connected" {
				fmt.Println("plugin: not connected — open Godot and enable the Godot MCP addon")
			} else {
				fmt.Println("plugin: connected")
			}
			return nil
		},
	}
}
