package version_test

import (
	"testing"

	"github.com/godot-mcp/godot-mcp/internal/version"
)

func TestParseSemVer(t *testing.T) {
	v, err := version.Parse("1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if v.Major != 1 || v.Minor != 0 || v.Patch != 0 {
		t.Fatalf("unexpected version: %+v", v)
	}
}

func TestCurrentVersion(t *testing.T) {
	if version.Version != "1.0.0" {
		t.Fatalf("expected 1.0.0, got %s", version.Version)
	}
}
