//go:build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/godot-mcp/godot-mcp/internal/bridge"
	"github.com/godot-mcp/godot-mcp/internal/godot"
	"github.com/godot-mcp/godot-mcp/internal/mcp"
	"github.com/godot-mcp/godot-mcp/internal/tools"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	godotConnectTimeout = 90 * time.Second
	godotStartupPause   = 3 * time.Second
)

// Env holds paths discovered for integration tests.
type Env struct {
	RepoRoot    string
	ProjectPath string
	GodotBin    string
}

// LoadEnv discovers repo root, demo project, and Godot binary.
func LoadEnv(t *testing.T) Env {
	t.Helper()

	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}

	project := filepath.Join(root, "examples", "demo")
	if _, err := os.Stat(filepath.Join(project, "project.godot")); err != nil {
		t.Fatalf("demo project missing: %v", err)
	}

	godotBin := os.Getenv("GODOT_BIN")
	if godotBin == "" {
		for _, name := range []string{"godot", "godot4", "Godot"} {
			if p, err := exec.LookPath(name); err == nil {
				godotBin = p
				break
			}
		}
	}
	if godotBin == "" {
		t.Skip("Godot not found — set GODOT_BIN or install Godot 4.3+")
	}

	return Env{
		RepoRoot:    root,
		ProjectPath: project,
		GodotBin:    godotBin,
	}
}

// Stack runs bridge, Godot editor, and MCP server against the demo project.
type Stack struct {
	Env      Env
	Addr     string
	Bridge   *bridge.WebSocket
	Client   *godot.Client
	MCPServer *mcp.Server
	MCPSession *mcpsdk.ClientSession
	cancel   context.CancelFunc
	godotCmd *exec.Cmd
}

// StartStack launches import, bridge, Godot editor, and an in-memory MCP client.
func StartStack(t *testing.T) *Stack {
	t.Helper()
	env := LoadEnv(t)

	if err := importProject(t, env); err != nil {
		if strings.Contains(err.Error(), "signal: killed") {
			t.Skipf("godot unavailable (killed by OS): %v", err)
		}
		t.Fatalf("godot import: %v", err)
	}

	addr, err := freeAddr()
	if err != nil {
		t.Fatalf("free addr: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	ws := bridge.NewWebSocket(addr)
	go func() {
		if err := ws.Start(ctx); err != nil && ctx.Err() == nil {
			t.Logf("bridge stopped: %v", err)
		}
	}()

	deadline := time.Now().Add(5 * time.Second)
	for {
		if resp, err := http.Get(fmt.Sprintf("http://%s/health", addr)); err == nil {
			resp.Body.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("bridge did not start")
		}
		time.Sleep(50 * time.Millisecond)
	}

	bridgeURL := fmt.Sprintf("ws://%s/ws", addr)
	cmd, output, exited := launchEditor(t, env, bridgeURL)
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})

	if err := waitForPlugin(ws, godotConnectTimeout, exited); err != nil {
		t.Fatalf("plugin connect: %v\n--- godot editor output ---\n%s", err, output.String())
	}

	client := godot.NewClient(ws)
	registry := tools.RegisterAll()
	mcpServer := mcp.New(mcp.Options{
		GodotClient: client,
		Registry:    registry,
	})

	clientTransport, serverTransport := mcpsdk.NewInMemoryTransports()
	go func() {
		_ = mcpServer.MCPServer().Run(ctx, serverTransport)
	}()

	mcpClient := mcpsdk.NewClient(&mcpsdk.Implementation{
		Name: "integration-test", Version: "1.0.0",
	}, nil)
	session, err := mcpClient.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("mcp connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	return &Stack{
		Env:        env,
		Addr:       addr,
		Bridge:     ws,
		Client:     client,
		MCPServer:  mcpServer,
		MCPSession: session,
		cancel:     cancel,
		godotCmd:   cmd,
	}
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func freeAddr() (string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr, nil
}

func importProject(t *testing.T, env Env) error {
	t.Helper()
	args := []string{"--headless", "--path", env.ProjectPath, "--import"}
	cmd := wrapDisplay(t, env.GodotBin, args...)
	cmd.Env = append(os.Environ(),
		"GODOT_MCP_ALLOW_DESTRUCTIVE=1",
		"GODOT_MCP_ALLOW_SCRIPT_EXEC=1",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

// syncBuffer is a goroutine-safe io.Writer for capturing subprocess output
// that's read concurrently (test failure message) while still being written to.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func launchEditor(t *testing.T, env Env, bridgeURL string) (*exec.Cmd, *syncBuffer, chan error) {
	t.Helper()
	args := []string{"--path", env.ProjectPath, "--editor"}
	cmd := wrapDisplay(t, env.GodotBin, args...)
	cmd.Env = append(os.Environ(),
		"GODOT_MCP_BRIDGE_URL="+bridgeURL,
		"GODOT_MCP_ALLOW_DESTRUCTIVE=1",
		"GODOT_MCP_ALLOW_SCRIPT_EXEC=1",
	)
	output := &syncBuffer{}
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Start(); err != nil {
		t.Fatalf("start godot editor: %v", err)
	}
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	time.Sleep(godotStartupPause)
	return cmd, output, exited
}

func wrapDisplay(t *testing.T, godotBin string, args ...string) *exec.Cmd {
	t.Helper()
	if runtime.GOOS == "linux" && os.Getenv("DISPLAY") == "" {
		if _, err := exec.LookPath("xvfb-run"); err == nil {
			all := append([]string{"-a", godotBin}, args...)
			return exec.Command("xvfb-run", all...)
		}
		t.Log("DISPLAY unset and xvfb-run missing — Godot editor may fail on Linux")
	}
	return exec.Command(godotBin, args...)
}

func waitForPlugin(ws *bridge.WebSocket, timeout time.Duration, exited chan error) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if ws.Connected() {
			info := ws.PluginInfo()
			if info != nil && info.ProjectName != "" {
				return nil
			}
		}
		select {
		case err := <-exited:
			return fmt.Errorf("godot editor exited before plugin connected: %v", err)
		case <-time.After(250 * time.Millisecond):
		}
	}
	return fmt.Errorf("godot plugin did not connect within %s", timeout)
}
