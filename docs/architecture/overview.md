# Architecture Overview

## Deployment model

Initial deployment is workstation-first:

```text
Engineer workstation or admin jumpbox
        ↓
esx9s TUI
        ↓
govmomi/govc/vSphere API
        ↓
standalone ESXi hosts
```

Running esx9s inside an ESXi-hosted VM is not the first target. That creates avoidable bootstrap and failure-domain problems.

## Main components

```text
cmd/esx9s        CLI/TUI entrypoint
internal/tui     terminal UI, views, keybindings
internal/esxi    ESXi/vSphere adapter
internal/config  host config and profiles
internal/inventory inventory cache/model
internal/actions plan/confirm/apply action model
internal/audit   append-only audit logging
internal/safety  confirmations, allowlists, destructive action guards
```

## Source of truth

ESXi/vSphere is the source of truth.

esx9s may cache data locally for speed, but cached state must not become authoritative unless explicitly designed and documented later.

## Backend strategy

Preferred backend:

- Go
- govmomi / govc-style API behavior
- vSphere API first

Possible constrained fallback:

- SSH for host-specific operations where API support is poor
- explicit allowlist only
- no raw command free-for-all in normal mode

## Action flow

```text
select object → choose action → build plan → confirm → apply → audit result
```

Destructive actions must never skip confirm/apply.

The audit event contract and example JSONL records are documented in
[Audit Log](audit-log.md).
