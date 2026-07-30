package filesystem

import (
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/permission"
	"github.com/godot-mcp/godot-mcp/internal/tools/toolutil"
)

func Register(r *mcp.Registry) {
	toolutil.RegisterWithParams[toolutil.OptionalPathParams](r, toolutil.Spec{
		Name: "file_list", Description: "List files in a directory.", RPC: "filesystem.list", Level: permission.LevelRead,
	}, func(p toolutil.OptionalPathParams) any {
		if p.Path == "" {
			p.Path = "res://"
		}
		return map[string]string{"path": p.Path}
	})
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "file_read", Description: "Read a text file.", RPC: "filesystem.read", Level: permission.LevelRead,
	}, nil)
	toolutil.RegisterWithParams[toolutil.FileWriteParams](r, toolutil.Spec{
		Name: "file_write", Description: "Write a text file.", RPC: "filesystem.write", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "file_create", Description: "Create an empty file.", RPC: "filesystem.create", Level: permission.LevelWrite,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "file_delete", Description: "Delete a file.", RPC: "filesystem.delete", Level: permission.LevelDestructive,
	}, nil)
	toolutil.RegisterWithParams[toolutil.PathParams](r, toolutil.Spec{
		Name: "folder_create", Description: "Create a folder.", RPC: "filesystem.mkdir", Level: permission.LevelWrite,
	}, nil)
}
