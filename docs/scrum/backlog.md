# Product Backlog

## Current focus

The current focus is `v0.1 - Operator Core`.

Do not pull future release scope into v0.1 unless it directly enables the first working operator loop.

## v0.1 - Operator Core

Goal: manage multiple standalone ESXi hosts safely from a first TUI.

Stories:

- Define product vision
- Record architecture decisions
- Scaffold Go project
- Define secure host config format
- Test connection to a host
- Build audit log foundation
- Build initial TUI shell
- List VMs across hosts
- Show VM detail view
- Implement VM power actions
- Implement snapshot create/delete

## v0.2 - Datastore Operator

Goal: make datastore work practical from the terminal.

Stories:

- List datastores per host
- Show usage/capacity
- Browse datastore paths
- Upload ISO
- Download file
- Register existing VMX

## v0.3 - VM Builder

Goal: create and configure VMs from the TUI.

Stories:

- Create VM from minimal config
- Configure CPU/RAM/disk
- Attach ISO
- Select port group/network
- Show creation plan before apply

## v0.4 - Host Operator

Goal: expose host-level day-2 operations.

Stories:

- Show host services
- Start/stop/restart safe host services
- Show host networking basics
- Enter/exit maintenance mode
- Show host events/tasks

## v0.5 - Fleet Operator

Goal: make multiple standalone hosts feel like one manageable fleet.

Stories:

- Offline move/copy VM between hosts
- Placement suggestions
- Resource pressure view
- Basic load-balancing recommendations
- Ops dashboard mode

## Parking lot

These are good ideas, but not v0.1:

- Kubernetes plugin
- Proxmox plugin
- OPNsense plugin
- TV dashboard mode
- mascot/brand polish
- advanced load balancing
- AI/advisor mode
- shared server mode
