package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "godot-mcp",
		Short: "MCP server for the Godot Engine",
	}

	root.AddCommand(newStartCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newVersionCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
