# Current Orchestration

Last reconciled: 2026-05-16

## Foundation Commit

Commit `50db545` is the foundation baseline:

- Mock TUI runs with `go run ./cmd/esx9s --mock`
- CI passes formatting, tests, vet, and build
- Config, audit, confirmation, mock provider, power, and snapshot foundations exist
- Real destructive ESXi operations are still intentionally deferred

## Closed Foundation Work

Issues closed as satisfied by `50db545`:

- #1, #3, #4, #5
- #6, #7, #8, #9, #10, #11, #12, #13
- #15, #17, #18, #20

## Open Work

Keep these open:

- #14: real provider-backed safe VM power operations and TUI action menu
- #16: real provider-backed snapshot operations and TUI action menu
- #19: GitHub Project board population, blocked until `gh auth refresh -s project`
- #21: versioning and release foundation epic
- #22: document versioning policy
- #23: add release checklist for first pre-release
- #24: add read-only govmomi connection implementation

## Next Execution Order

1. Finish project-board wiring after granting the `project` scope.
2. Complete versioning policy and release checklist.
3. Tag the first pre-release from a clean, green `main`.
4. Implement read-only govmomi connection and inventory.
5. Only then continue real VM power and snapshot actions.

## Board Blocker

The GitHub project exists at `esx9s Kanban`, but `gh project item-add` fails without the write-oriented `project` scope.

Run:

```sh
gh auth refresh -s project
```

Then add the open issues to the board and move closed foundation work to Done.
