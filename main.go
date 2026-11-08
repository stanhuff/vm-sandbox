package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type sandboxInfo struct {
	ID string `json:"id"`
}

func run(ctx context.Context, args []string, in io.ReadCloser, out io.WriteCloser) error {
	switch {
	case len(args) == 3 && args[0] == "serve":
		return serve(ctx, args[1], args[2])
	case len(args) == 2 && args[0] == "mcp":
		return serveMCP(ctx, args[1], in, out)
	default:
		return errors.New("usage: vm-sandbox serve CONFIG_DIR SOCKET_DIR | vm-sandbox mcp SOCKET")
	}
}

func serveMCP(ctx context.Context, socket string, in io.ReadCloser, out io.WriteCloser) error {
	server := mcp.NewServer(&mcp.Implementation{Name: "vm-sandbox", Version: "0.1.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sandbox_info",
		Description: "Identify the sandbox associated with this connection.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, sandboxInfo, error) {
		info, err := readSandbox(ctx, socket)
		return nil, info, err
	})
	err := server.Run(ctx, &mcp.IOTransport{Reader: in, Writer: out})
	if ctx.Err() != nil {
		return nil
	}
	return err
}

func readSandbox(ctx context.Context, socket string) (sandboxInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "unix", socket)
	if err != nil {
		return sandboxInfo{}, err
	}
	defer conn.Close()
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()

	var info sandboxInfo
	err = json.NewDecoder(io.LimitReader(conn, 1024)).Decode(&info)
	if err == nil && info.ID == "" {
		err = errors.New("missing sandbox ID")
	}
	return info, err
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
