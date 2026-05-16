# GitHub project setup

Issue #19 completed the repo-side metadata setup:

- Labels created: `type:task`, `type:story`, `type:spike`, `release:v0.1`, `area:github`.
- Milestone created: `v0.1`.
- `release:v0.1` and the `v0.1` milestone were applied to issues #6, #7, #8, #9, #10, #12, #14, and #16.
- Type labels were applied to issues #11, #13, #15, #17, #18, and #19.
- `area:github` was applied to issue #19.

The `esx9s Kanban` project board exists at:

https://github.com/users/nielslindor/projects/4

Direct project mutation was blocked because the active `gh` token is missing the `project` scope. To finish populating the board, refresh the token and add the open issues:

```sh
gh auth refresh -s project

for issue in 1 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do
  gh project item-add 4 \
    --owner nielslindor \
    --url "https://github.com/nielslindor/esx9s/issues/$issue"
done
```

Recommended initial status for all added items is `Backlog`.
