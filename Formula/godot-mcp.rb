class GodotMcp < Formula
  desc "Model Context Protocol server for Godot Engine 4.x"
  homepage "https://github.com/triforge0/godot-mcp"
  license "MIT"
  head "https://github.com/triforge0/godot-mcp.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = "-s -w"
    system "go", "build", *std_go_args(ldflags:), "./cmd/godot-mcp"
    pkgshare.install "plugin"
    pkgshare.install "examples/demo"
  end

  def caveats
    <<~EOS
      The Godot editor plugin is installed at:
        #{pkgshare}/plugin/addons/godot_mcp

      Copy or symlink it into your Godot project:
        ln -s #{pkgshare}/plugin/addons/godot_mcp /path/to/project/addons/godot_mcp

      Then enable "Godot MCP" under Project → Project Settings → Plugins.
    EOS
  end

  test do
    assert_match "godot-mcp", shell_output("#{bin}/godot-mcp version")
    assert_path_exists pkgshare/"plugin/addons/godot_mcp/plugin.cfg"
  end
end
