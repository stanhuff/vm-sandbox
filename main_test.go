package main

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSandboxInfoOverStdio(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config := t.TempDir()
	sockets, err := os.MkdirTemp("", "vms-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(sockets)
	for _, id := range []string{"alpha", "beta"} {
		if err := os.WriteFile(filepath.Join(config, id+".toml"), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}

	done := make(chan error, 1)
	go func() { done <- run(ctx, []string{"serve", config, sockets}, nil, nil) }()
	defer func() {
		cancel()
		if err := <-done; err != nil {
			t.Error(err)
		}
	}()

	for _, id := range []string{"alpha", "beta"} {
		socket := filepath.Join(sockets, id, "control.sock")
		for {
			if _, err := os.Stat(socket); err == nil {
				break
			}
			select {
			case <-ctx.Done():
				t.Fatal("socket did not appear:", socket)
			case <-time.After(5 * time.Millisecond):
			}
		}
		t.Run(id, func(t *testing.T) {
			serverIn, clientOut := io.Pipe()
			clientIn, serverOut := io.Pipe()
			adapterDone := make(chan error, 1)
			go func() {
				adapterDone <- run(ctx, []string{"mcp", socket}, serverIn, serverOut)
				serverOut.Close()
			}()
			defer func() {
				clientOut.Close()
				clientIn.Close()
				if err := <-adapterDone; err != nil {
					t.Error(err)
				}
			}()

			client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil)
			session, err := client.Connect(ctx, &mcp.IOTransport{Reader: clientIn, Writer: clientOut}, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer session.Close()
			result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "sandbox_info"})
			if err != nil {
				t.Fatal(err)
			}
			info, ok := result.StructuredContent.(map[string]any)
			if result.IsError || !ok || info["id"] != id {
				t.Fatalf("expected sandbox %s, got %+v", id, result)
			}
		})
	}
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"unknown"}, {"serve"}, {"mcp"}} {
		if err := run(context.Background(), args, nil, nil); err == nil || !strings.Contains(err.Error(), "usage") {
			t.Errorf("run(%v) = %v, want usage error", args, err)
		}
	}
}

func TestInvalidSandboxConfigs(t *testing.T) {
	for _, tc := range []struct {
		name, body string
	}{
		{"bad name.toml", ""},
		{"alpha.toml", "unknown = true"},
		{"alpha.toml", "socket_gid = -1"},
		{"alpha.toml", "socket_gid = ["},
	} {
		t.Run(tc.name+"/"+tc.body, func(t *testing.T) {
			config := t.TempDir()
			if err := os.WriteFile(filepath.Join(config, tc.name), []byte(tc.body), 0600); err != nil {
				t.Fatal(err)
			}
			if err := serve(context.Background(), config, t.TempDir()); err == nil {
				t.Fatal("expected invalid config to fail")
			}
		})
	}
	if err := serve(context.Background(), filepath.Join(t.TempDir(), "missing"), t.TempDir()); err == nil {
		t.Fatal("expected missing config directory to fail")
	}
}

func TestSocketPermissionsAndCollision(t *testing.T) {
	root, err := os.MkdirTemp("", "vms-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	gid := os.Getgid()
	listener, err := listenSandbox(root, "alpha", sandboxConfig{SocketGID: &gid})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	path := filepath.Join(root, "alpha", "control.sock")
	for name, mode := range map[string]os.FileMode{path: 0660, filepath.Dir(path): 0710} {
		info, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != mode {
			t.Errorf("%s mode = %o, want %o", name, info.Mode().Perm(), mode)
		}
	}
	if second, err := listenSandbox(root, "alpha", sandboxConfig{}); err == nil {
		second.Close()
		t.Fatal("second listener opened an existing socket")
	}
	if err := os.Symlink(filepath.Dir(path), filepath.Join(root, "beta")); err != nil {
		t.Fatal(err)
	}
	if second, err := listenSandbox(root, "beta", sandboxConfig{}); err == nil {
		second.Close()
		t.Fatal("listener followed a symlink")
	}
}

func TestReadSandboxFailures(t *testing.T) {
	root, err := os.MkdirTemp("", "vms-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	path := filepath.Join(root, "control.sock")
	if _, err := readSandbox(context.Background(), path); err == nil {
		t.Fatal("missing socket unexpectedly succeeded")
	}
	listener, err := net.Listen("unix", path)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			defer conn.Close()
			io.WriteString(conn, "{}")
		}
	}()
	if _, err := readSandbox(context.Background(), path); err == nil {
		t.Fatal("missing sandbox ID unexpectedly succeeded")
	}
}
