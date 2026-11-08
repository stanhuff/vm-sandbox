# vm-sandbox

The first development slice connects an MCP client to a configured sandbox and reports its identity.

In VS Code, select **Dev Containers: Reopen in Container**. The image includes Go and the C toolchain used by the race detector. From the container terminal, run `go test -race ./...` or build with `go build -o vm-sandbox .`.

Create an operator-controlled config directory and socket root. Place one TOML file per sandbox in the config directory. For example, `appa-dev.toml` defines the sandbox `appa-dev` and can be empty. IDs use lowercase ASCII letters, digits, and hyphens, start with a letter or digit, and have at most 32 characters.

Run `vm-sandbox serve CONFIG_DIR SOCKET_DIR` on the host. Bind-mount only the desired `SOCKET_DIR/appa-dev` directory read-only into the devcontainer. Configure its MCP client to launch `vm-sandbox mcp /path/to/mount/control.sock` over stdio and call `sandbox_info`.

By default, only the daemon's host UID can connect. To grant access to a group, put `socket_gid = 1234` in the sandbox TOML and use the container client's host-mapped GID. Run the daemon under an account authorized to assign that group. Keep both host directories under operator control.
