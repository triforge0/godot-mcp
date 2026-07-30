package domain

// GodotObject is the generic representation of any Godot engine object.
// Tools and reflection APIs use this instead of hard-coded node types.
type GodotObject struct {
	Class      string         `json:"class"`
	Path       string         `json:"path,omitempty"`
	Name       string         `json:"name,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	Children   []GodotObject  `json:"children,omitempty"`
}

type PropertyValue struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    any    `json:"value"`
	Type     string `json:"type,omitempty"`
}

type MethodInfo struct {
	Name        string `json:"name"`
	ReturnType  string `json:"return_type,omitempty"`
	ArgumentCount int  `json:"argument_count,omitempty"`
}

type ClassInfo struct {
	Class      string         `json:"class"`
	Properties []PropertyMeta `json:"properties,omitempty"`
	Methods    []MethodInfo   `json:"methods,omitempty"`
}

type PropertyMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type RPCResult map[string]any
