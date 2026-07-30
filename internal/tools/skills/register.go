package skills

import (
	"context"

	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
)

func Register(r *mcp.Registry) {
	r.Register(mcp.ToolDefinition{
		Name:        "skill_create_player",
		Description: "Create a 2D player (CharacterBody2D + CollisionShape2D + Sprite2D).",
		RPCMethod:   "skill.create_player",
		Permission:  permission.LevelWrite,
		Handler:     createPlayer,
	})
	r.Register(mcp.ToolDefinition{
		Name:        "skill_create_scene",
		Description: "Create and open a new basic scene.",
		RPCMethod:   "skill.create_scene",
		Permission:  permission.LevelWrite,
		Handler:     createScene,
	})
	r.Register(mcp.ToolDefinition{
		Name:        "skill_optimize_project",
		Description: "Analyze project for common issues and optimization hints.",
		RPCMethod:   "skill.optimize_project",
		Permission:  permission.LevelRead,
		Handler: func(ctx context.Context, client *godot.Client, _ any) (any, error) {
			return client.Call(ctx, "skill.optimize_project", nil, nil)
		},
	})
	r.Register(mcp.ToolDefinition{
		Name:        "skill_analyze_error",
		Description: "Analyze recent errors and suggest fixes.",
		RPCMethod:   "skill.analyze_error",
		Permission:  permission.LevelRead,
		Handler: func(ctx context.Context, client *godot.Client, params any) (any, error) {
			return client.Call(ctx, "skill.analyze_error", params, nil)
		},
	})
}

func createPlayer(ctx context.Context, client *godot.Client, params any) (any, error) {
	p, _ := params.(map[string]any)
	if p == nil {
		p = map[string]any{}
	}
	return client.CallMap(ctx, "skill.create_player", p)
}

func createScene(ctx context.Context, client *godot.Client, params any) (any, error) {
	p, _ := params.(map[string]any)
	return client.CallMap(ctx, "skill.create_scene", p)
}
