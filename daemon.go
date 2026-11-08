package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pelletier/go-toml/v2"
)

var sandboxID = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)

type sandboxConfig struct {
	SocketGID *int `toml:"socket_gid"`
}

func loadSandboxes(dir string) (map[string]sandboxConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sandboxes := make(map[string]sandboxConfig)
	for _, entry := range entries {
		id, found := strings.CutSuffix(entry.Name(), ".toml")
		if !found {
			continue
		}
		if !sandboxID.MatchString(id) || !entry.Type().IsRegular() {
			return nil, fmt.Errorf("invalid sandbox configuration file %q", entry.Name())
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var config sandboxConfig
		if err := toml.NewDecoder(bytes.NewReader(data)).DisallowUnknownFields().Decode(&config); err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		if config.SocketGID != nil && (*config.SocketGID < 0 || uint64(*config.SocketGID) >= 1<<32-1) {
			return nil, fmt.Errorf("%s: socket_gid is outside the valid range", entry.Name())
		}
		sandboxes[id] = config
	}
	return sandboxes, nil
}

func listenSandbox(root, id string, config sandboxConfig) (*net.UnixListener, error) {
	dir := filepath.Join(root, id)
	if err := os.Mkdir(dir, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}
	path := filepath.Join(dir, "control.sock")
	if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("socket path %s is already in use or inaccessible", path)
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return nil, err
	}

	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		return nil, err
	}
	if err := setSocketAccess(dir, path, config.SocketGID); err != nil {
		listener.Close()
		return nil, err
	}
	return listener, nil
}

func setSocketAccess(dir, path string, gid *int) error {
	if gid == nil {
		return os.Chmod(path, 0600)
	}
	if err := os.Chown(dir, -1, *gid); err != nil {
		return err
	}
	if err := os.Chown(path, -1, *gid); err != nil {
		return err
	}
	if err := os.Chmod(path, 0660); err != nil {
		return err
	}
	return os.Chmod(dir, 0710)
}

func serve(ctx context.Context, configDir, socketRoot string) error {
	configs, err := loadSandboxes(configDir)
	if err != nil {
		return err
	}
	listeners := make(map[string]*net.UnixListener)
	defer func() {
		for _, listener := range listeners {
			listener.Close()
		}
	}()
	for id, config := range configs {
		listener, err := listenSandbox(socketRoot, id, config)
		if err != nil {
			return fmt.Errorf("sandbox %s: %w", id, err)
		}
		listeners[id] = listener
	}

	var workers sync.WaitGroup
	failures := make(chan error, len(listeners))
	for id, listener := range listeners {
		workers.Go(func() {
			failures <- respond(listener, id)
		})
	}
	select {
	case <-ctx.Done():
	case err = <-failures:
	}
	for _, listener := range listeners {
		listener.Close()
	}
	workers.Wait()
	return err
}

func respond(listener *net.UnixListener, id string) error {
	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			return err
		}
		conn.SetWriteDeadline(time.Now().Add(time.Second))
		json.NewEncoder(conn).Encode(sandboxInfo{ID: id})
		conn.Close()
	}
}
