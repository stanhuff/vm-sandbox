# Sandbox access

## REQ-0001 — Independent sandbox configuration

An operator can define an empty sandbox using a dedicated configuration file. The operator can define multiple sandboxes on one host.

## REQ-0002 — Sandbox identity over MCP

A devcontainer with access to a sandbox socket can connect over MCP and query the sandbox's ID. The daemon determines the identity from the socket used for the connection.
