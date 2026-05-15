# Project Operating Model

## Tooling decision

Use GitHub as the source of truth for esx9s.

Reason:

- The repo already lives on GitHub
- Code, issues, pull requests, reviews, CI, releases, documentation, and the project board can stay together
- GitHub Issues and Projects are enough for the current project size
- Adding Asana now would create duplicate task state
- GitLab is not needed unless the repository moves there later

## Asana decision

Do not use Asana for active esx9s execution yet.

Asana may become useful later for high-level business planning, partnership work, content, launch planning, or non-code company operations. It should not mirror GitHub issues.

Rule:

> GitHub owns product and engineering execution. Asana, if added later, owns business/process work only.

## Working model

Use a lightweight Scrum/Kanban hybrid:

- Product vision lives in `docs/product/`
- Architecture decisions live in `docs/adr/`
- Project rules live in `docs/scrum/`
- Work items live in GitHub Issues
- The active board lives in GitHub Projects: `esx9s Kanban`
- Code changes happen through pull requests
- Releases are tracked by milestones or roadmap items

## Board columns

Use the active GitHub Project columns:

1. Backlog
2. Ready
3. In progress
4. In review
5. Done

Optional later column:

- Parked

If Parked does not exist, keep parked ideas in Backlog and mark them clearly as parked in the issue body or `docs/scrum/drift-keeping.md`.

## Issue types

Use title prefixes until GitHub issue types/labels are fully configured:

- `[Epic]`
- `[Story]`
- `[Task]`
- `[Spike]`
- `[Bug]`
- `[Security]`
- `[Docs]`

## Priorities

Use issue body text or Project custom fields:

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

- a GitHub issue in the `esx9s Kanban` project, or
- intentionally parked in `docs/scrum/drift-keeping.md`

No invisible backlog in chat.
