# esx9s

A k9s-inspired terminal operator console for standalone ESXi fleets.

esx9s is not a k9s fork and not an open-source vCenter clone. The product goal is narrower and sharper: give homelabbers, small environments, and ops-focused engineers a lightweight, secure, beautiful TUI for managing multiple standalone ESXi hosts without requiring vCenter.

## Quickstart

Prerequisites:

- Go 1.24 or newer
- Make

Run the mock TUI:

```sh
make run
```

Useful development commands:

```sh
make test
make fmt
make build
```

To use a specific local Go toolchain, set `GO`:

```sh
GO=/Users/mijnbeen/.cache/esx9s-tools/go/bin/go make test
```

Print the scaffold version:

```sh
make build
./bin/esx9s --version
```

Run mock mode directly:

```sh
go run ./cmd/esx9s --mock
```

Probe an ESXi SDK endpoint:

```sh
go run ./cmd/esx9s connect-test --endpoint https://HOST/sdk
```

## Product positioning

**One sentence:** esx9s is a lightweight terminal operator for standalone ESXi fleets.

**Operator focus:**

- Manage multiple standalone ESXi hosts
- Unified VM inventory
- Create and configure VMs
- Power actions
- Snapshots
- Datastore browse/upload/download
- ISO upload
- Register/unregister VMs
- Offline move/copy between hosts
- Host networking basics
- Host services
- Resource usage
- Event/task view
- SSH/API health checks
- Safe action audit log
- Placement and load-balancing suggestions later
- Beautiful ops dashboard mode later

## Boundaries

esx9s should not pretend to replace every vCenter feature. Some capabilities are real vCenter constructs. The correct boundary is:

> k9s for standalone ESXi fleets, not a full vCenter clone.

## First release target

`v0.1 - Operator Core`

- multi-host config
- connection tests
- unified VM inventory
- VM detail view
- power actions
- snapshot create/delete
- append-only audit log
- secure-by-default operator flow

## Development principles

- Workstation-first app
- Go-first implementation
- govmomi/govc backend first
- SSH fallback only where useful and constrained
- No stored passwords
- No hidden destructive actions
- Plan/confirm/apply for dangerous actions
- Every operator action is auditable
- Plugin-ready later, but ESXi-first now
