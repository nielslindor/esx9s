# Sprint 0 - Project Foundation

## Goal

Create enough structure that esx9s can be built without drifting into a generic infrastructure platform.

## Sprint outcome

At the end of Sprint 0, the project should have:

- clear product vision
- accepted architecture decisions
- Go project scaffold
- CI build
- initial TUI shell
- secure host config shape
- connection test spike
- audit log foundation

## Sprint rules

- No Kubernetes plugin work
- No mascot implementation work
- No VM creation work
- No datastore upload/download work
- No load-balancing implementation
- No TV dashboard implementation

Those are parking-lot items until the operator core exists.

## Definition of done for Sprint 0

- README explains the product clearly
- docs/product and docs/architecture exist
- ADRs for Go, govmomi/govc, and workstation-first exist
- backlog exists
- first runnable binary exists
- `make test` or equivalent exists
- at least one host connection can be tested or stubbed
- audit log format is documented or implemented

## Primary risk

Scope expansion.

If a task does not directly help build the first operator loop, it should not enter Sprint 0.

First operator loop:

```text
configure host → test connection → list VMs → select VM → perform safe action → audit result
```
