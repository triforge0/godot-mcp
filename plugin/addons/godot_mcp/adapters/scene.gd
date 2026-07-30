class_name MCPSceneAdapter
extends RefCounted

const B = preload("res://addons/godot_mcp/adapter_base.gd")


static func handles(method: String) -> bool:
	return method.begins_with("scene.")


static func call(method: String, params: Variant) -> Variant:
	match method:
		"scene.list":
			return {"scenes": _list_scenes("res://")}
		"scene.current":
			var root := B.edited_root()
			if root == null:
				return {"scene": null}
			return {"scene": root.scene_file_path, "root": root.name}
		"scene.open":
			var path := str(B.param(params, "path", ""))
			if path.is_empty():
				return B.err("path required")
			var err := EditorInterface.open_scene_from_path(path)
			return {"opened": err == OK, "path": path} if err == OK else B.err("open failed")
		"scene.save":
			if B.edited_root() == null:
				return B.err("no scene open")
			var err := EditorInterface.save_scene()
			return {"saved": err == OK}
		"scene.create":
			var root := Node2D.new()
			root.name = str(B.param(params, "name", "NewScene"))
			EditorInterface.edit_node(root)
			return {"created": true, "name": root.name}
		"scene.reload":
			var root := B.edited_root()
			if root == null or root.scene_file_path.is_empty():
				return B.err("no saved scene")
			EditorInterface.open_scene_from_path(root.scene_file_path)
			return {"reloaded": true}
		"scene.close":
			EditorInterface.close_scene()
			return {"closed": true}
	return B.err("unsupported: %s" % method)


static func _list_scenes(path: String) -> Array:
	var out: Array = []
	var dir := DirAccess.open(path)
	if dir == null:
		return out
	dir.list_dir_begin()
	var f := dir.get_next()
	while f != "":
		var full := path.path_join(f)
		if dir.current_is_dir() and not f.begins_with("."):
			out.append_array(_list_scenes(full))
		elif f.ends_with(".tscn") or f.ends_with(".scn"):
			out.append(full)
		f = dir.get_next()
	dir.list_dir_end()
	return out
