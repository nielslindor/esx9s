# Safety Model

esx9s is not read-only, so safety is a product feature, not a later hardening task.

## Safety goals

- prevent accidental destructive actions
- avoid credential leakage
- make operator intent visible
- make action history auditable
- keep ESXi/vSphere as source of truth

## Rules

1. No stored plaintext passwords.
2. No secrets in logs.
3. Destructive actions require confirmation.
4. Actions produce audit events.
5. API-backed actions are preferred over shell-backed actions.
6. SSH fallback must be explicit and allowlisted.
7. Raw shell mode is not normal product behavior.

## Action model

```text
select object → choose action → build plan → show target → confirm → apply → audit result
```

## Confirmation levels

### Level 0: Read-only

Examples:

- list hosts
- list VMs
- show datastore usage
- show events/tasks

No confirmation needed.

### Level 1: Reversible or low-risk action

Examples:

- power on VM
- graceful shutdown
- create snapshot

Requires visible target confirmation.

### Level 2: Disruptive action

Examples:

- hard power off
- reset VM
- delete snapshot
- unregister VM

Requires typed confirmation.

### Level 3: Destructive action

Examples:

- delete VM files
- delete datastore file
- destructive host network change

Out of scope for v0.1 unless explicitly designed.

## Audit event fields

Minimum JSONL fields:

```json
{
  "timestamp": "2026-05-15T12:00:00Z",
  "operator": "local-user",
  "target_host": "esxi05",
  "target_type": "vm",
  "target_id": "vm-123",
  "target_name": "opnsense",
  "action": "power_off",
  "plan_id": "uuid",
  "result": "success",
  "error": null
}
```
