# esx9s Product Vision

## What this is

esx9s is a k9s-inspired terminal operator console for ESXi fleets with vCenter-like day-2 functionality.

It gives operators a fast, keyboard-driven, secure way to manage multiple ESXi hosts with the operational coverage people normally reach for vCenter to get. Standalone ESXi is the first backend; optional vCenter-backed behavior can be added later where vCenter is the correct API boundary.

## North star

Oneshot the whole product: keep decomposing, implementing, testing, documenting, and reconciling project state until esx9s is a usable k9s-for-ESXi control plane with vCenter-like operational coverage.

## What this is not

- Not a k9s fork
- Not a vCenter server clone or UI clone
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
- Perform vCenter-like day-2 operations through safe terminal workflows

## Design feel

- k9s-like speed and keyboard feel
- vCenter-like operational coverage where ESXi or vCenter APIs can honestly support it
- compact binary / lightweight install
- safe enough for real ops
- pretty enough to keep on a wall-mounted terminal or TV

## Strategic line

The win is not recreating VMware vCenter as a server product.

The win is:

> I can manage ESXi fleets from a beautiful, safe, fast TUI with the day-2 functionality I normally associate with vCenter.
