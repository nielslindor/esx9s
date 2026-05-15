# GitHub Projects Board Usage

## Decision

Use GitHub Projects as the operational board for esx9s.

Asana is intentionally not part of the active workflow yet. The project is still small enough that duplicating work between GitHub and Asana would create drift and stale task state.

## Board

Project: `esx9s Kanban`

Recommended current views:

- Backlog
- Priority board
- Team items
- Roadmap
- My items

## Columns

Use the current board columns as follows.

### Backlog

Ideas or valid work that has not been refined enough to start.

Use for:

- feature ideas
- parked scope
- rough stories
- issues without acceptance criteria

Do not pull directly from Backlog into In progress. First refine it and move it to Ready.

### Ready

Work that can be picked up without more thinking.

A Ready issue has:

- clear outcome
- acceptance criteria
- non-goals
- priority or milestone intent
- no unresolved architecture question blocking it

### In progress

Actively being worked on.

Limit this hard. For solo work, keep this at 1 or 2 items. The board may allow more, but esx9s should not.

### In review

A PR exists, or the work needs review/validation before Done.

Use for:

- code review
- docs review
- architecture review
- manual testing against ESXi

### Done

The issue meets the Definition of Done and is closed.

## Issue hierarchy

Use parent/child structure where possible:

- Epics represent releases or major product areas
- Stories represent operator-visible outcomes
- Tasks represent implementation steps
- Spikes represent time-boxed research
- Bugs represent broken expected behavior

For v0.1, use one parent epic:

- `[Epic] v0.1 - Operator Core`

Child issues:

- Scaffold Go project
- Define host config format
- Build ESXi connection test
- Build audit log foundation
- Build initial TUI shell

## Dependencies

Use dependencies when one issue blocks another.

Initial dependency chain:

```text
Scaffold Go project
  -> Define host config format
  -> Build ESXi connection test
  -> Build initial TUI shell

Build audit log foundation
  -> VM power actions later
  -> snapshot actions later
```

If GitHub dependency fields are not enabled, write dependencies explicitly in the issue body using:

```text
Blocked by: #123
Blocks: #456
```

## Roadmap view

Use the Roadmap view for milestones/releases, not every small task.

Recommended roadmap items:

- v0.1 - Operator Core
- v0.2 - Datastore Operator
- v0.3 - VM Builder
- v0.4 - Host Operator
- v0.5 - Fleet Operator
- v1.0 - Stable Small Fleet Release

## Priority board

Use Priority board for current execution pressure.

Recommended priority meaning:

- P0: blocks project, security, or release integrity
- P1: needed for current milestone
- P2: valuable but later

## My items

Use My items as the daily execution surface.

Rule:

> If it is not on My items or intentionally parked, it is not today's work.

## Drift-keeping rule

When a good idea appears, do not immediately add it to In progress.

Choose one:

1. Add as Backlog issue
2. Add to `docs/scrum/drift-keeping.md`
3. Attach as a note to the relevant epic

Default for new ideas:

> Backlog first, Ready only after acceptance criteria exist.

## Sprint 0 board setup

Move these to Ready after checking the acceptance criteria:

- `[Task] Scaffold Go project`
- `[Story] Define host config format`
- `[Spike] Build ESXi connection test`
- `[Story] Build audit log foundation`
- `[Story] Build initial TUI shell`

Keep these in Done or close when verified:

- Define esx9s product vision
- Document workstation-first architecture
- Document Go plus govmomi/govc backend decision

Keep project board/admin tasks in Backlog or Ready depending on whether manual setup remains.
