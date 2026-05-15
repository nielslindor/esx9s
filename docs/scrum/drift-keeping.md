# Drift Keeping

Drift keeping is the project habit of catching scope expansion before it turns into unfinished complexity.

## Product anchor

esx9s is a k9s-inspired terminal operator console for standalone ESXi fleets.

The sharp product boundary:

> k9s for standalone ESXi fleets, not an open-source vCenter clone.

## Current target

The current target is `v0.1 - Operator Core`.

v0.1 proves:

- multi-host config works
- ESXi connection testing works
- unified VM inventory works
- safe VM power actions work
- snapshot create/delete works
- audit logging works
- the TUI shell feels plausible

## Drift signals

Pause when work shifts toward:

- Kubernetes plugins before ESXi core works
- Proxmox/OPNsense/Veeam before ESXi core works
- mascot/art direction before the first TUI shell works
- load balancing before inventory and actions work
- AI/advisor suggestions before basic telemetry works
- full vCenter replacement claims
- raw SSH command playgrounds
- enterprise suite packaging before the product exists

## Drift response

When drift appears, ask:

1. Does this help `v0.1 - Operator Core` ship?
2. Can it become a later issue instead?
3. What is the smallest next product proof?

Default answer:

> Park it. Ship the operator core.

## Parked but valid ideas

These are good ideas, just not Sprint 0:

- datastore browser/upload/download
- ISO upload
- VM creation/configuration
- offline move/copy between hosts
- placement/load-balancing suggestions
- wallboard/TV mode
- Kubernetes plugin
- mascot / pixel raccoon operator
- broader StackDeck-style suite
