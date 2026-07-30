# Testing

## Unit tests

```bash
make test
# or
go test ./...
```

## Integration tests (real Godot editor)

Integration tests launch the Godot 4 editor with the demo project, connect the MCP plugin over WebSocket, and exercise tools through both RPC and the MCP SDK.

### Prerequisites

- Godot **4.3+** editor binary
- Linux CI/local headless: `xvfb-run` (Linux only, when `DISPLAY` is unset)

### Install Godot (optional helper)

```bash
bash scripts/install-godot.sh
# binary default: /usr/local/bin/godot
# custom dir: GODOT_INSTALL_DIR=$HOME/.local/bin bash scripts/install-godot.sh
```

### Run

```bash
make test-integration
# or
GODOT_BIN=/path/to/godot go test -tags=integration -timeout 10m ./internal/integration/...
```

Tests skip automatically when Godot is not installed.

### CI

The `integration-godot` job in [.github/workflows/ci.yml](../.github/workflows/ci.yml) installs Godot 4.3 on Ubuntu with `xvfb` and runs integration tests on every push/PR.

## End-to-end mock tests

Mock WebSocket plugin tests (no Godot required):

```bash
go test ./internal/e2e/...
```
