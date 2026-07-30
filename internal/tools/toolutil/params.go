package toolutil

// Shared MCP parameter types used across tool families.

type PathParams struct {
	Path string `json:"path" jsonschema:"res:// path or node path"`
}

type OptionalPathParams struct {
	Path string `json:"path,omitempty" jsonschema:"optional res:// path"`
}

type SceneCreateParams struct {
	Name string `json:"name,omitempty" jsonschema:"root node name, default NewScene"`
	Path string `json:"path,omitempty" jsonschema:"optional res:// path to save the new scene"`
}

type NodeCreateParams struct {
	ParentPath string `json:"parent_path" jsonschema:"path to the parent node"`
	Type       string `json:"type" jsonschema:"Godot class name e.g. CharacterBody2D"`
	Name       string `json:"name" jsonschema:"name for the new node"`
}

type NodeRenameParams struct {
	Path    string `json:"path" jsonschema:"path to the node"`
	NewName string `json:"new_name" jsonschema:"new node name"`
}

type NodeMoveParams struct {
	Path          string `json:"path" jsonschema:"path to the node to move"`
	NewParentPath string `json:"new_parent_path" jsonschema:"path to the new parent node"`
	Index         int    `json:"index,omitempty" jsonschema:"optional child index under new parent"`
}

type PropertyGetParams struct {
	Path     string `json:"path" jsonschema:"path to the node"`
	Property string `json:"property" jsonschema:"property name"`
}

type PropertySetParams struct {
	Path     string `json:"path" jsonschema:"path to the node"`
	Property string `json:"property" jsonschema:"property name"`
	Value    any    `json:"value" jsonschema:"property value"`
}

type FileWriteParams struct {
	Path    string `json:"path" jsonschema:"res:// file path"`
	Content string `json:"content" jsonschema:"file content"`
}

type ScriptSourceParams struct {
	Path   string `json:"path" jsonschema:"res:// script path"`
	Source string `json:"source,omitempty" jsonschema:"GDScript source code"`
}

type ScriptAttachParams struct {
	Path       string `json:"path" jsonschema:"path to the node"`
	ScriptPath string `json:"script_path" jsonschema:"res:// path to the script"`
}

type ResourceCreateParams struct {
	Path string `json:"path" jsonschema:"res:// output path"`
	Type string `json:"type" jsonschema:"Resource class name"`
}

type ResourceListParams struct {
	Path string `json:"path,omitempty" jsonschema:"directory to list, default res://"`
}

type ClassInspectParams struct {
	Class string `json:"class" jsonschema:"Godot class name e.g. Sprite2D"`
}

type SkillCreateSceneParams struct {
	Name string `json:"name,omitempty" jsonschema:"scene name, default NewScene"`
	Path string `json:"path,omitempty" jsonschema:"res:// save path"`
}

type SkillCreatePlayerParams struct {
	ParentPath string `json:"parent_path,omitempty" jsonschema:"parent node path, defaults to scene root"`
}

type SkillAnalyzeErrorParams struct {
	Limit int `json:"limit,omitempty" jsonschema:"max errors to analyze, default 5"`
}
