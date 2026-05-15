# ADR 0001: Use Go as primary implementation language

## Status

Accepted for initial implementation.

## Context

esx9s should feel lightweight, compact, and easy to distribute. It should not feel like a heavy Python toolchain or a pile of scripts.

The project needs:

- compact binaries
- strong terminal/TUI ecosystem
- good concurrency for multi-host operations
- good VMware API ecosystem
- easy cross-platform builds

## Decision

Use Go as the primary implementation language.

Likely TUI libraries:

- Bubble Tea
- Bubbles
- Lip Gloss

Likely VMware backend:

- govmomi
- govc-inspired behavior

## Consequences

Positive:

- compact release artifacts
- good fit for CLI/TUI tools
- good ecosystem for vSphere through govmomi
- easy install story later

Negative:

- less rapid prototyping than Python
- TUI polish requires discipline

## Explicit rejection

pyVmomi is rejected for the core product direction because it gives a heavier Python/sludge feeling and does not match the compact k9s-like product identity.
