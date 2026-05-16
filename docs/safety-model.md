# Safety Model

Safety is part of the first slice.

- Do not store passwords in config.
- Do not print passwords, tokens, session cookies, or authorization headers.
- Destructive actions require explicit confirmation.
- Operator actions write append-only JSONL audit events.
- Real ESXi writes stay behind provider interfaces and must pass through plan, confirm, apply, and audit.

The current implementation includes typed confirmation helpers, an audit JSONL writer, and mock-only power/snapshot action paths. Real provider writes are intentionally deferred.
