# Definition of Done

A work item is done only when all relevant points below are true.

## Code

- Builds locally
- Tests pass where tests exist
- No secrets or credentials committed
- Errors are handled with useful operator messages
- Action paths are auditable where applicable

## Product

- The implemented behavior matches the issue acceptance criteria
- Dangerous actions require explicit confirmation
- UX is keyboard-friendly
- Failure modes are visible, not silent

## Docs

- New behavior is documented where needed
- Architecture changes get an ADR if they affect long-term design
- New config fields are reflected in example config

## Security

- No stored plaintext passwords
- No secrets in logs
- No unaudited destructive actions
- Any SSH usage is allowlisted or explicitly marked as break-glass

## Review

- PR description explains what changed
- PR includes screenshots/terminal captures for TUI changes when useful
- Scope creep is rejected or moved to a separate issue
