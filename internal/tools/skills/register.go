package skills

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterWithParams[toolutil.SkillCreatePlayerParams](r, toolutil.Spec{
		Name: "skill_create_player", Description: "Create a 2D player (CharacterBody2D + collision + sprite).",
		RPC: "skill.create_player", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterWithParams[toolutil.SkillCreateSceneParams](r, toolutil.Spec{
		Name: "skill_create_scene", Description: "Create and open a new basic scene.",
		RPC: "skill.create_scene", Level: permission.LevelWrite,
	}, nil)

	toolutil.RegisterNoParams(r, toolutil.Spec{
		Name: "skill_optimize_project", Description: "Analyze project for common issues and optimization hints.",
		RPC: "skill.optimize_project", Level: permission.LevelRead,
	})

	toolutil.RegisterWithParams[toolutil.SkillAnalyzeErrorParams](r, toolutil.Spec{
		Name: "skill_analyze_error", Description: "Analyze recent errors and suggest fixes.",
		RPC: "skill.analyze_error", Level: permission.LevelRead,
	}, func(p toolutil.SkillAnalyzeErrorParams) any {
		if p.Limit <= 0 {
			p.Limit = 5
		}
		return p
	})
}
