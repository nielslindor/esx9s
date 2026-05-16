# ADR 0002: Use govmomi/govc as backend foundation

## Status

Accepted for initial implementation.

## Context

esx9s needs to manage ESXi/vSphere objects safely and consistently.

Direct SSH commands are tempting, but they increase parsing fragility, permission ambiguity, and command safety risk. The app should feel like an operator console, not a remote shell wrapper.

## Decision

Use govmomi as the primary Go library and govc behavior as the operational reference model.

SSH may exist later as a constrained, allowlisted fallback for host-specific tasks, but it must not be the core abstraction.

## Consequences

Positive:

- better API semantics
- less brittle command parsing
- clearer object model
- strong fit with Go
- easier to support both standalone ESXi and vCenter later

Negative:

- some host operations may still need SSH or esxcli-equivalent workflows
- API auth and session handling must be implemented carefully

## Guardrail

No normal-mode raw SSH command runner.

If a break-glass shell mode ever exists, it must be explicit, visibly dangerous, and audited.
