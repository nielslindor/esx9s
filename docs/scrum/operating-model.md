# Project Operating Model

## Tooling decision

Use GitHub as the source of truth for esx9s.

Reason:

- The repo already lives on GitHub
- Code, issues, pull requests, reviews, CI, releases, and documentation can stay together
- GitHub Issues are enough for the current project size
- Adding Asana now would create duplicate task state
- GitLab is not needed unless the repository moves there later

## Working model

Use a lightweight Scrum/Kanban hybrid:

- Product vision lives in `docs/product/`
- Architecture decisions live in `docs/adr/`
- Project rules live in `docs/scrum/`
- Work items live in GitHub Issues
- Code changes happen through pull requests
- Releases are tracked by milestones

## Board columns

Recommended GitHub Project columns:

1. Inbox
2. Ready
3. In Progress
4. Review
5. Done
6. Parked

## Issue types

Use title prefixes until labels are created:

- `[Epic]`
- `[Story]`
- `[Task]`
- `[Spike]`
- `[Bug]`
- `[Security]`
- `[Docs]`

## Priorities

Use issue body text until labels are created:

- P0: blocks the project or security
- P1: needed for current milestone
- P2: useful but later

## Sprint rhythm

For solo/founder-style work:

- Sprint length: 1 week
- Sprint planning: pick 3 to 5 issues max
- Daily check: choose one next action, not a full replanning session
- Review: close what is actually done, park drift

## Current sprint

Current sprint: `Sprint 0 - Project Foundation`

Goal: make the repo structured enough that implementation can begin without scope drift.

## Source of truth rule

If it is real work, it must be either:

- a GitHub issue, or
- intentionally parked in `docs/scrum/drift-keeping.md`

No invisible backlog in chat.
