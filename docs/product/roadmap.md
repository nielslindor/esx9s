# Roadmap

## Product goal

Oneshot esx9s into the full k9s-for-ESXi terminal control plane with vCenter-like day-2 functionality.

The roadmap is not a suggestion to stop at the next slice. It is the durable path Codex should keep executing across context compression, commits, issues, tests, and releases.

## v0.1 - Operator Core

Goal: prove that esx9s can safely control multiple standalone ESXi hosts from a TUI.

Scope:

- multi-host config
- connection test
- unified VM inventory
- VM detail view
- power actions
- snapshot create/delete
- append-only audit log
- first TUI shell

Out of scope:

- VM creation
- datastore transfer
- VM movement
- advisor suggestions
- Kubernetes plugins
- wallboard/TV mode

## v0.2 - Datastore Operator

Goal: make datastore work practical from the terminal.

Scope:

- datastore list
- datastore usage
- datastore browser
- upload ISO
- download file
- register VMX
- safe datastore action preview and audit

## v0.3 - VM Builder

Goal: cover vCenter-like VM lifecycle work from the TUI.

Scope:

- create VM
- configure CPU/RAM/disk
- attach ISO
- basic network selection
- safe plan screen before creation
- register/unregister VM
- clone/copy/offline move where technically possible

## v0.4 - Host Operator

Goal: expose useful host-level operations, networking basics, and events/tasks.

Scope:

- host services
- basic networking view
- maintenance mode
- host health checks
- events/tasks
- allowlisted safe host service actions

## v0.5 - Fleet Operator

Goal: make multiple standalone hosts feel like one manageable fleet.

Scope:

- offline move/copy VM between hosts
- placement suggestions
- basic load-balancing recommendations
- fleet health overview
- cross-host search and filtering

## v1.0 - Stable Small Fleet Release

Goal: reliable release for homelab and small fleet users.

Scope:

- install packages
- config migration
- stable docs
- release pipeline
- safety guarantees documented
- common failure modes handled
- capability matrix for standalone ESXi versus optional vCenter-backed behavior
