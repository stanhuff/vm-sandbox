# Glossary

- **Operator** — The person who administers the VM host, defines sandboxes, and grants clients access to them.
- **VM host** — The machine that runs the virtual machines and the vm-sandbox daemon.
- **Sandbox** — A named boundary for access to VM resources and for quota accounting. A sandbox can exist without any VMs or sandbox-owned images. Multiple clients can share a sandbox.
- **Sandbox configuration** — The operator-managed configuration file that defines one sandbox. Its settings include resource quotas and the information needed to grant access.
- **Quota** — A limit on resource use by a sandbox. Requests from clients sharing a sandbox count against the same limits.
- **Virtual machine (VM)** — A guest computer managed on the VM host. Each managed VM belongs to one sandbox.
- **Image** — A reusable disk template from which a VM can be created. A base image is supplied by the operator. A derived image is published from a VM and belongs to its sandbox.
- **Daemon** — The host-side service that loads sandbox configurations, enforces access and quotas, and manages resources on behalf of clients.
- **Devcontainer** — A development environment running in a container. It can host an AI agent and the tools the agent uses.
- **AI agent** — The software that performs development tasks and calls tools to interact with vm-sandbox.
- **MCP client** — The component used by an agent to discover and call tools via the Model Context Protocol.
- **MCP adapter** — The process that presents vm-sandbox tools as an MCP server to a client and forwards requests to the daemon.
- **Sandbox socket** — The host-side Unix socket through which an adapter connects to a particular sandbox.
