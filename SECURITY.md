# Security Policy

esx9s is an infrastructure operator tool. Security is core product behavior, not a later add-on.

## Supported versions

No stable release exists yet.

## Security principles

- No stored plaintext passwords
- No secrets in logs
- No unaudited destructive actions
- No normal-mode raw SSH command runner
- Dangerous operations require confirmation
- ESXi/vSphere remains source of truth
- Local cache must not become hidden authority
- Every operator action should be auditable

## Credentials

Initial direction:

- Prefer session-based auth or OS/keychain-backed secrets later
- Avoid writing passwords to config files
- Config examples must not contain real credentials
- SSH fallback, if implemented, must use allowlisted operations

## Audit logging

Operator actions should produce append-only audit events including:

- timestamp
- local user where available
- target host
- target object
- action
- planned operation
- result
- error, if any

Audit logs must not include:

- passwords
- tokens
- private keys
- session cookies

## Reporting vulnerabilities

Until a formal disclosure process exists, open a private security advisory if GitHub security advisories are enabled. Otherwise contact the maintainer directly.
