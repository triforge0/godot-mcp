class_name MCPRuntimeAdapter
extends RefCounted

const B = preload("res://addons/godot_mcp/adapter_base.gd")


static func handles(method: String) -> bool:
	return method.begins_with("runtime.")


static func call(method: String, params: Variant) -> Variant:
	match method:
		"runtime.run":
			EditorInterface.play_main_scene()
			return {"running": true}
		"runtime.stop":
			EditorInterface.stop_playing_scene()
			return {"stopped": true}
		"runtime.pause":
			if EditorInterface.is_playing_scene():
				EditorInterface.set_pause_playing_scene(true)
			return {"paused": true}
		"runtime.resume":
			if EditorInterface.is_playing_scene():
				EditorInterface.set_pause_playing_scene(false)
			return {"resumed": true}
		"runtime.status":
			return {
				"playing": EditorInterface.is_playing_scene(),
				"fps": Engine.get_frames_per_second(),
			}
		"runtime.restart":
			EditorInterface.stop_playing_scene()
			EditorInterface.play_main_scene()
			return {"restarted": true}
	return B.err("unsupported: %s" % method)
