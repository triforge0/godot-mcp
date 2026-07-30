class_name MCPScreenshotAdapter
extends RefCounted

const B = preload("res://addons/godot_mcp/adapter_base.gd")


static func handles(method: String) -> bool:
	return method.begins_with("screenshot.")


static func call(method: String, params: Variant) -> Variant:
	match method:
		"screenshot.capture", "screenshot.viewport":
			var img := _capture()
			if img == null or img.is_empty():
				return B.err("capture failed")
			var png := img.save_png_to_buffer()
			return {
				"format": "png",
				"data": Marshalls.raw_to_base64(png),
				"width": img.get_width(),
				"height": img.get_height(),
			}
	return B.err("unsupported: %s" % method)


static func _capture() -> Image:
	if EditorInterface.is_playing_scene():
		var playing := EditorInterface.get_playing_scene()
		if playing:
			var tex := playing.get_viewport().get_texture()
			if tex:
				return tex.get_image()
	for getter in [EditorInterface.get_editor_viewport_2d, EditorInterface.get_editor_viewport_3d]:
		var vp = getter.call()
		if vp:
			var tex := vp.get_texture()
			if tex:
				return tex.get_image()
	return null
