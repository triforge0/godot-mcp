package version

import (
	"fmt"
	"strconv"
	"strings"
)

const Name = "godot-mcp"

// Version is the released version. Release builds override it via
// -ldflags "-X github.com/godot-mcp/godot-mcp/internal/version.Version=X.Y.Z".
var Version = "1.0.0"

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func Parse(v string) (SemVer, error) {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	parts := strings.Split(v, "-")
	core := strings.Split(parts[0], ".")
	if len(core) != 3 {
		return SemVer{}, fmt.Errorf("invalid semver: %s", v)
	}
	major, err := strconv.Atoi(core[0])
	if err != nil {
		return SemVer{}, err
	}
	minor, err := strconv.Atoi(core[1])
	if err != nil {
		return SemVer{}, err
	}
	patch, err := strconv.Atoi(core[2])
	if err != nil {
		return SemVer{}, err
	}
	return SemVer{Major: major, Minor: minor, Patch: patch}, nil
}

func (s SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func Current() SemVer {
	v, _ := Parse(Version)
	return v
}
