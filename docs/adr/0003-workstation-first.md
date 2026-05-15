# ADR 0003: Workstation-first deployment

## Status

Accepted for initial implementation.

## Context

esx9s is an operator tool for controlling ESXi hosts. It could theoretically run inside a VM on ESXi, but that creates a failure-domain and bootstrap problem: if that host, datastore, VM, or network breaks, the control plane may be unavailable when needed most.

## Decision

The first serious version runs on an engineer workstation or hardened admin jumpbox.

```text
engineer workstation / admin jumpbox
        ↓
esx9s
        ↓
ESXi/vSphere API
        ↓
standalone ESXi hosts
```

## Consequences

Positive:

- easier development
- safer failure model
- easier to kill/upgrade/debug
- aligns with per-engineer credentials
- avoids installing open-source control software inside the ESXi estate

Negative:

- each engineer may need local install/config
- shared dashboard mode comes later

## Future option

A server/jumpbox mode may be added later, but not before the workstation TUI is useful.
