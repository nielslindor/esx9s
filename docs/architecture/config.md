# Config Format

esx9s uses a workstation-local YAML file to describe the standalone ESXi hosts an operator can connect to.

The v0.1 config format is intentionally small: it identifies hosts, connection endpoints, operator usernames, authentication lookup methods, and TLS verification behavior. It must not contain passwords or reusable secrets.

By default, the loader looks for `config.yaml` under the current user's `esx9s` config directory as reported by the operating system.

## Example

See `configs/esx9s.example.yaml` for a complete multi-host example.

```yaml
version: 1

hosts:
  - name: lab-esxi-01
    address: esxi01.lab.example.com
    endpoint: https://esxi01.lab.example.com/sdk
    username: root
    auth:
      method: prompt
    tls:
      insecure_skip_verify: false
```

## Top-level fields

| Field | Required | Description |
| --- | --- | --- |
| `version` | Yes | Config schema version. v0.1 uses `1`. |
| `hosts` | Yes | Ordered list of ESXi hosts available to the operator. |

## Host fields

| Field | Required | Description |
| --- | --- | --- |
| `name` | Yes | Stable local identifier shown in the TUI and audit events. Must be unique within the file. |
| `address` | Yes | Hostname or IP address for display, matching, and connection diagnostics. |
| `endpoint` | Yes | vSphere SDK endpoint, normally `https://<address>/sdk`. |
| `username` | Yes | Default operator username for this host. This is not secret. |
| `auth` | Yes | Authentication method and lookup metadata. Must not include secret values. |
| `tls` | Yes | TLS verification settings. |

## Authentication methods

The config describes how esx9s should obtain credentials. It never stores the credential values themselves.

### `prompt`

Prompt the operator at connection time.

```yaml
auth:
  method: prompt
```

Use this as the safest default and for shared examples.

### `keychain`

Read credentials from the operating system credential store.

```yaml
auth:
  method: keychain
  service: esx9s
  account: edge-esxi-01
```

The exact backend is platform-specific. On macOS this maps naturally to Keychain; on other platforms the loader can later map the same method to a secure local credential store.

Additional authentication methods can be added later if they only reference secure external lookup mechanisms and do not introduce password fields into the YAML schema.

## TLS verification

TLS verification is enabled by default:

```yaml
tls:
  insecure_skip_verify: false
```

`insecure_skip_verify: true` disables certificate verification and makes man-in-the-middle attacks easier. It is acceptable only for isolated labs, first-contact bootstrap, or temporary self-signed certificate situations. Production-like environments should install trusted host certificates and keep this value `false`.

The TUI should visibly warn when a host uses `insecure_skip_verify: true`, especially before applying actions.

## Security rules

1. Do not add `password`, `pass`, `secret`, `token`, session cookie, private key, or similar fields to the config schema.
2. Do not log resolved credential values.
3. Do not include credentials in audit events.
4. Prefer per-operator credentials over shared root credentials.
5. Keep file permissions restrictive because hostnames, usernames, and environment variable names can still reveal operational details.

## Validation expectations

The loader rejects configs that:

- omit a host `name`
- omit a host `address`
- include `password`, `pass`, `secret`, or `token` field names

The loader may support additional non-secret metadata later, but secret material must remain outside the YAML file.
