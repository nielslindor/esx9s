# esx9s Product Vision

## What this is

esx9s is a k9s-inspired terminal operator console for standalone ESXi fleets.

It gives operators a fast, keyboard-driven, secure way to manage multiple ESXi hosts without needing to deploy vCenter for small fleets, labs, edge environments, and ops-heavy infrastructure work.

## What this is not

- Not a k9s fork
- Not a vCenter clone
- Not a full cluster scheduler
- Not a Kubernetes plugin platform yet
- Not a Terraform replacement
- Not a toy read-only dashboard

## Target user

Primary:

- Homelab operators
- Infra/network engineers
- Small environments running multiple standalone ESXi hosts
- Edge/site operators who need practical VM and host control

Secondary:

- MSP-style operators
- Internal lab teams
- Engineers who want vCenter-like day-2 operations without vCenter overhead

## Product promise

From one beautiful TUI, an operator can:

- See all hosts and VMs
- Act on VMs safely
- Work with datastores
- Upload ISOs
- Register existing VMs
- Move/copy VMs between hosts where technically possible
- Understand host health and resource pressure
- Review events/tasks
- Keep an audit trail of actions

## Design feel

- k9s-like speed and keyboard feel
- vCenter-like operational coverage for standalone hosts
- compact binary / lightweight install
- safe enough for real ops
- pretty enough to keep on a wall-mounted terminal or TV

## Strategic line

The win is not recreating VMware vCenter.

The win is:

> I can manage three to ten standalone ESXi hosts from a beautiful, safe, fast TUI without spinning up vCenter.
