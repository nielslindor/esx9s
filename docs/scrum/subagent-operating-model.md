# Subagent Operating Model

The goal is to oneshot esx9s into a k9s-for-ESXi control plane with vCenter-like day-2 functionality. Subagents are how Codex keeps that large goal moving without losing state or creating merge chaos.

## Prime Rule

Use subagents as a coordinated engineering squad, not as a crowd.

Every subagent must have:

- one GitHub issue or clearly named research question
- a bounded output
- an explicit write set, or read-only status
- a final report with files changed, tests run, blockers, and issue/project updates needed

## Startup Pattern

At the start of a substantial session:

1. Read `docs/scrum/current-orchestration.md`.
2. Inspect open issues and `esx9s Kanban`.
3. Choose the critical path issue.
4. Decide what the main Codex thread should do locally.
5. Spawn only sidecar work that can run in parallel without blocking the next local step.

Do not delegate the immediate blocker if the main thread cannot move without it.

Use the project board column and issue body blockers as truth. Labels can be stale or broad; if a `status:ready` label disagrees with the board column or a `Blocked by:` line, fix the board and issue state before spawning work.

## Subagent Brief Packet

Every subagent prompt should include:

- issue URL and number
- current base commit
- project board column
- blockers and dependencies
- acceptance criteria and non-goals
- allowed write set, or explicit read-only status
- relevant docs to inspect
- tests or commands expected before final report
- reminder that secrets, lab host details, and credentials must not be written to files, logs, issue comments, or final output

## Standard Squad Shapes

### Planning / Roadmap Squad

Use when creating or reshaping issues, milestones, capability matrices, or release scope.

- Project orchestration worker: issue/board state, sequencing, dependencies
- Product scope worker: roadmap, capability matrix, acceptance criteria
- Merger: updates docs/issues/project after both reports

Good for: #22, #23, #25, #37, #38, #39, #40, #41, #42, #43.

### Provider Squad

Use when adding real ESXi or future vCenter-backed behavior.

- Credential worker: credential source API, redaction tests, docs
- Provider worker: govmomi client/inventory/actions in `internal/provider`
- Integration-test worker: opt-in harness and documentation
- Safety/QA worker: secret scan, failure modes, audit expectations
- Merger: integrates provider/config/TUI seams and runs full verification

Good for: #26, #27, #28, #29, #30, #36.

### TUI Workflow Squad

Use when building operator screens or action flows.

- TUI worker: Bubble Tea model/views/key handling
- Action worker: plan/confirm/apply domain logic
- Audit/safety worker: audit path and destructive-action gates
- QA worker: model tests, noninteractive render tests, keyboard smoke notes
- Merger: owns `cmd/esx9s`, integration tests, and final user command

Good for: #31, #32, #33, #34, #14, #16.

### Release Squad

Use before tagging or publishing.

- Release docs worker: versioning policy, release checklist, notes
- CI worker: workflow, secret/lab-detail guard, build matrix
- QA worker: local and GitHub CI verification
- Merger: tag/release only after green checks and clean status

Good for: #21, #22, #23, #25, #35.

## Current Parallelization Guidance

Parallelize now:

- #22 and #23 as a release-docs lane
- #35 as a CI/security lane
- #26 as the credential/provider prerequisite lane
- #43 as the capability-matrix lane

Do not spawn these yet, even if labels say ready:

- #25 until #22 and #23 are done
- #27 until #26 is done
- #28 until #27 is done
- #29 until #27 has stable env names and client boundaries
- #14 or #16 until all write-action safety gates are complete

## Ownership Boundaries

Use disjoint write sets whenever possible.

- CLI: `cmd/esx9s/*`
- TUI: `internal/tui/*`, `internal/app/*`
- Domain/actions: `internal/domain/*`, `internal/actions/*`, `internal/confirm/*`
- Provider: `internal/provider/*`
- Config: `internal/config/*`, `configs/*`, config docs
- Audit/safety: `internal/audit/*`, safety docs
- CI/release: `.github/workflows/*`, `Makefile`, release docs
- Product/project docs: `docs/product/*`, `docs/scrum/*`, issue/project metadata

High-conflict files:

- `go.mod` and `go.sum`
- `cmd/esx9s/main.go`
- `internal/domain/models.go`
- `internal/tui/model.go`
- `docs/scrum/current-orchestration.md`

Only the merger should make broad edits to high-conflict files after workers finish.

Shared contract files require extra care:

- `internal/domain/models.go`
- `internal/provider/provider.go`
- `internal/confirm/confirm.go`
- audit event shapes in `internal/audit` and `internal/domain`

If a worker needs to change a shared contract, prefer landing that contract first, then fan out dependent workers.

## Merger Role

Every multi-agent coding push needs a merger.

The merger owns:

- reviewing worker outputs
- resolving conflicts and duplicated abstractions
- running `gofmt`, `go test ./...`, `go vet ./...`, and relevant smoke commands
- checking for committed secrets or lab details
- updating issue comments/project states
- creating blocker issues when a worker cannot be merged cleanly
- producing the final commit or PR-ready state

If a merge blocker appears, do not leave it in chat. Create a GitHub issue with the failing command, files involved, and exact unblock criteria.

## QA and Safety Gates

Use a dedicated QA/safety subagent before real ESXi writes or releases.

Required checks:

- `go test ./...`
- `go test -race ./...` before risky provider/action changes when practical
- `go vet ./...`
- `go run ./cmd/esx9s --mock`
- secret/lab-detail scan
- audit path verification for action flows
- confirmation flow verification for destructive actions
- opt-in integration tests only when explicitly enabled
- release checks before tagging: clean worktree, green GitHub CI, clean secret scan, accurate release notes

Real ESXi credentials may only be used by the designated provider/integration worker for the issue at hand. They must be supplied by environment variable or prompt and must never appear in files, test fixtures, logs, issue comments, shell snippets, or release notes.

## Issue and Board Discipline

Subagents do not just write code.

For each issue they touch, they must report:

- issue number
- status: done, partial, blocked, or follow-up needed
- files changed
- tests/commands run
- suggested issue comment
- board transition needed

The main thread or merger applies issue/project updates. Do not let issue state drift behind code state.

## Compression Handoff

Before the main thread risks losing context:

1. Stop or close unneeded subagents.
2. Commit or clearly leave a clean worktree state.
3. Update `docs/scrum/current-orchestration.md`.
4. Update GitHub issues and project board.
5. Write the next exact execution step into the relevant issue.

Future Codex sessions should be able to resume from GitHub plus repo docs without reading chat history.
