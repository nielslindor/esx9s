# GitHub Actions Automation Templates

The active workflow files should live in `.github/workflows/`.

This folder contains copy-ready templates for repo automation. They are stored as documentation first so they can be reviewed before activation.

## Recommended activation order

1. `bootstrap-labels.yml` - manually run once to create labels.
2. `ci.yml` - enables repo smoke checks and Go checks once `go.mod` exists.
3. `issue-triage.yml` - auto-labels issues and flags missing acceptance criteria.
4. `pr-hygiene.yml` - requires PRs to reference issues.
5. `project-auto-add.yml` - optional; requires a `GH_PROJECT_TOKEN` secret for GitHub Projects v2 access.

## Why templates first

GitHub Actions can execute code in the repository context. Review before activation.

## Project board automation split

Use GitHub Projects built-in workflows for board status movement where possible.
Use GitHub Actions for repo automation, CI, label hygiene, PR hygiene, and optional project sync.
