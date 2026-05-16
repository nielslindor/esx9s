# Scope

esx9s is a k9s-like terminal control plane for multiple standalone ESXi hosts.

## Current Slice

- Mock-driven Go TUI
- Hosts, VMs, Datastores, and Tasks/Events views
- Local config shape and validation
- Append-only local audit log
- Typed confirmation helpers
- Provider interface with mock implementation
- Read-only ESXi SDK reachability probe

## Not Now

- vCenter replacement claims
- Kubernetes or plugin marketplace work
- Web dashboard
- Raw SSH command runner
- Real destructive ESXi actions without confirmation and audit
