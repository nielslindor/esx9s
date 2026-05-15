# Product Principles

## 1. Operator-first

esx9s is for active operators, not passive viewers. The TUI should support real action, but action must be deliberate and auditable.

## 2. Safe by default

Dangerous operations require a clear plan, visible target, and explicit confirmation.

## 3. Fast enough to live in

The app should feel instant for normal inventory browsing and common actions.

## 4. Standalone ESXi first

Do not dilute the project with Kubernetes, Proxmox, OPNsense, or generic plugin ambitions before the ESXi core works.

## 5. Workstation-first

The first serious version runs on the engineer workstation or admin jumpbox, not inside ESXi as a VM.

## 6. vCenter-like usefulness, not vCenter mimicry

Copy the useful operator outcomes, not the full platform surface.

## 7. No hidden state magic

Local cache is allowed. Hidden authoritative state is not. ESXi/vSphere remains source of truth.

## 8. Audit everything

Every action should produce an append-only audit event with timestamp, user, target, action, and result.

## 9. Pretty matters

The interface should be clear, readable, and visually satisfying. This is an ops tool people should want to keep open.
