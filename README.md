# vm-sandbox

The first development slice connects an MCP client to a configured sandbox and reports its identity.

In VS Code, select **Dev Containers: Reopen in Container**. The image includes Go and the C toolchain used by the race detector. From the container terminal, run `go test -race ./...` or build with `go build -o vm-sandbox .`.

Create an operator-controlled config directory and socket root. Place one TOML file per sandbox in the config directory. For example, `appa-dev.toml` defines the sandbox `appa-dev` and can be empty. IDs use lowercase ASCII letters, digits, and hyphens, start with a letter or digit, and have at most 32 characters.

Run `vm-sandbox serve CONFIG_DIR SOCKET_DIR` on the host. Bind-mount only the desired `SOCKET_DIR/appa-dev` directory read-only into the devcontainer. Configure its MCP client to launch `vm-sandbox mcp /path/to/mount/control.sock` over stdio and call `sandbox_info`.

By default, only the daemon's host UID can connect. To grant access to a group, put `socket_gid = 1234` in the sandbox TOML and use the container client's host-mapped GID. Run the daemon under an account authorized to assign that group. Keep both host directories under operator control.

## Ubuntu package

Download the `vm-sandbox-deb` artifact from a successful GitHub Actions run, or run
`dpkg-buildpackage -b -us -uc` inside the VS Code devcontainer. The local build
writes the package one directory above the checkout.

Install the extracted package on the Ubuntu VM host with
`sudo apt install ./vm-sandbox_0.1.0_amd64.deb` (use the filename for your architecture).
The package installs both commands and starts the `vm-sandbox` systemd service.
Add your sandbox TOML files under `/etc/vm-sandbox/sandboxes.d` and run
`sudo systemctl restart vm-sandbox`. For a non-root client, set `socket_gid` to
the client's host-mapped group. Mount only the matching sandbox directory under
`/run/vm-sandbox`. Recreate the client container after a service restart so the
bind mount picks up the newly created runtime directory.
