# Contributing to esx9s

esx9s is early-stage. The most important contribution rule is to keep scope tight.

## Project boundary

esx9s is:

- a k9s-inspired TUI operator console
- for standalone ESXi fleets
- focused on practical day-2 operations
- secure and auditable by default

esx9s is not:

- a full vCenter clone
- a Kubernetes platform
- a generic infrastructure suite yet
- a collection of random SSH scripts

## Work style

Use issues for work. Every issue should have:

- outcome
- acceptance criteria
- explicit non-goals
- safety impact if applicable

## Pull requests

PRs should be small and focused. A PR should ideally close one issue.

Include screenshots or terminal captures for TUI changes when useful.

## Architecture changes

Long-term technical decisions should get an ADR in `docs/adr/`.

## Security-sensitive work

Anything touching auth, credentials, destructive actions, SSH, audit logs, or host control needs extra care and must be reviewed against `docs/scrum/definition-of-done.md`.
