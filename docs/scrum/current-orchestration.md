# Current Orchestration

Last reconciled: 2026-05-16

## Durable Context Rule

The GitHub issues and `esx9s Kanban` project are the source of truth. Do not rely on chat history for sequencing, blocked state, or scope boundaries. Context compression will happen; the project has to remember the mission.

The mission is not merely the next issue. The mission is to oneshot the entire k9s-for-ESXi control plane with vCenter-like day-2 functionality.

When a Codex session starts:

1. Read this file.
2. Inspect open GitHub issues.
3. Inspect the project board.
4. Pick the highest-priority Ready issue that is not blocked.
5. Update issues/project state as work changes.

## Foundation Baseline

Commit `50db545` is the mock foundation baseline:

- `go run ./cmd/esx9s --mock` renders the mock TUI.
- CI passes formatting, tests, vet, and build.
- Config, audit, confirmation, mock provider, mock power, and mock snapshot foundations exist.
- Real destructive ESXi operations are intentionally deferred.

Commit `e36a227` added this orchestration handoff.

## Board State

Project: `esx9s Kanban`

The board now contains:

- Closed foundation issues in `Done`
- Current release/provider/action issues in `Ready` or `Backlog`
- Priority and size fields for the active lane

If the board and this file disagree, reconcile the board first and then update this file.

## Closed Foundation Work

Issues closed as satisfied by `50db545`:

- #1, #3, #4, #5
- #6, #7, #8, #9, #10, #11, #12, #13
- #15, #17, #18, #20
- #19 after project-board population

## Open Epics

- #21: Versioning and release foundation
- #24: Real ESXi read-only provider
- #37: Oneshot esx9s vCenter-like terminal control plane
- #38: v0.2 Datastore and ISO operator
- #39: v0.3 VM lifecycle operator
- #40: v0.4 Host, networking, tasks, and events operator
- #41: v0.5 Fleet operations and placement
- #42: v1.0 Stable vCenter-like small fleet release

## Open Work Queue

Release/versioning:

- #22: Document versioning policy
- #23: Add release checklist for first pre-release
- #25: Create first pre-release tag

Provider/read-only ESXi:

- #26: Add secure credential input layer
- #27: Implement govmomi login/logout client
- #28: Implement read-only govmomi inventory provider
- #29: Add opt-in ESXi integration test harness
- #30: Wire config mode into the TUI
- #36: Add safe real-ESXi quickstart

TUI/operator workflow:

- #31: Add VM detail view
- #32: Add TUI VM action menus
- #33: Add action preview and confirmation screen
- #34: Wire default audit log path into action flows

Safety/CI:

- #35: Add CI secret and lab-detail guard
- #43: Build vCenter functionality capability matrix

Real write actions:

- #14: Implement safe VM power actions
- #16: Implement snapshot create/delete flow

Long-range epics:

- #37: Oneshot full product
- #38: Datastore and ISO operator
- #39: VM lifecycle operator
- #40: Host/networking/tasks/events operator
- #41: Fleet operations and placement
- #42: Stable vCenter-like small fleet release

## Execution Order

Do this next:

1. #22 versioning policy
2. #23 release checklist
3. #43 vCenter functionality capability matrix
4. #35 CI secret and lab-detail guard
5. #25 first pre-release tag
6. #26 credential input layer
7. #27 govmomi login/logout
8. #29 opt-in integration harness
9. #28 read-only inventory provider
10. #30 config-mode TUI
11. #31 VM detail view
12. #34 default audit path
13. #33 action preview/confirmation screen
14. #32 TUI VM action menus
15. #14 real safe VM power actions
16. #16 real snapshot actions
17. #36 safe real-ESXi quickstart, updated with whatever reality was learned
18. Continue through #38, #39, #40, #41, and #42 until the full product is complete.

## Safety Gates

Do not implement real ESXi write actions until all of these are done:

- #26 secure credential input
- #28 read-only provider inventory
- #33 action preview and confirmation screen
- #34 default audit log wiring
- #35 CI secret/lab-detail guard

No real ESXi hostnames, IPs, usernames, passwords, tokens, session cookies, or Authorization headers may be committed.

## Scope Guard

Still not now:

- Kubernetes
- Web UI
- AI remediation
- HA/vMotion/DRS replacement claims
- Raw SSH command runner

Allowed as the north star:

- vCenter-like day-2 functionality from a terminal UI
- optional vCenter provider/backend later where standalone ESXi cannot honestly provide a capability
- explicit capability matrix to prevent false claims

Park useful side ideas as issues with `phase:later` and `drift-control`, then return to the current queue.
