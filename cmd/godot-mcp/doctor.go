package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/godot-mcp/godot-mcp/internal/config"
	"github.com/godot-mcp/godot-mcp/internal/tools"
	"github.com/godot-mcp/godot-mcp/internal/version"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose Godot MCP connectivity and configuration",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := config.Load()
			registry := tools.RegisterAll()

			fmt.Printf("godot-mcp %s\n", version.Version)
			fmt.Printf("tools registered: %d\n", registry.Len())
			fmt.Printf("bridge address: %s\n", cfg.BridgeAddr)

			client := &http.Client{Timeout: 2 * time.Second}
			resp, err := client.Get(fmt.Sprintf("http://%s/health", cfg.BridgeAddr))
			if err != nil {
				fmt.Println("\nStatus: NOT READY")
				fmt.Println("  bridge: not reachable — ensure `godot-mcp start` runs via your MCP client")
				fmt.Println("  plugin: enable Godot MCP in Project → Project Settings → Plugins")
				return nil
			}
			defer resp.Body.Close()

			var health map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
				return fmt.Errorf("decode health: %w", err)
			}

			fmt.Printf("\nbridge: reachable (%s)\n", health["status"])
			switch health["status"] {
			case "connected":
				fmt.Println("Status: READY")
				printIfSet("plugin version", health["plugin_version"])
				printIfSet("godot version", health["godot_version"])
				printIfSet("project", health["project_name"])
			default:
				fmt.Println("Status: WAITING FOR PLUGIN")
				fmt.Println("  plugin: not connected — open Godot and enable the Godot MCP addon")
			}

			fmt.Println("\nPermissions:")
			fmt.Println("  Destructive tools and script_execute prompt in the Godot editor.")
			fmt.Println("  Dev bypass (skip dialogs):")
			fmt.Println("    GODOT_MCP_ALLOW_DESTRUCTIVE=1")
			fmt.Println("    GODOT_MCP_ALLOW_SCRIPT_EXEC=1")
			return nil
		},
	}
}

func printIfSet(label, value string) {
	if value != "" {
		fmt.Printf("%s: %s\n", label, value)
	}
}
