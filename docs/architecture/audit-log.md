# Audit Log

esx9s writes operator actions as append-only JSONL records. Each line is one
complete audit event so the file can be streamed, tailed, archived, or imported
without parsing multi-line records.

## Required fields

Each audit event must include:

- `timestamp`: UTC timestamp for the attempted action.
- `operator`: local operator identity when available.
- `target_host`: configured host name or address for the action target.
- `target_type`: object kind, such as `vm`, `host`, `snapshot`, or `datastore`.
- `target_id`: stable provider identifier for the target object.
- `target_name`: human-friendly name when available.
- `action`: requested operation.
- `plan_id`: identifier for the confirmed plan that produced the action.
- `result`: outcome such as `success`, `failed`, or `denied`.
- `error`: sanitized error text, or `null` when no error occurred.

## Example records

Successful VM power action:

```json
{"timestamp":"2026-05-15T12:00:00Z","operator":"local-user","target_host":"esxi05","target_type":"vm","target_id":"vm-123","target_name":"opnsense","action":"power_off","plan_id":"plan-123","result":"success","error":null}
```

Failed snapshot action:

```json
{"timestamp":"2026-05-15T12:03:21Z","operator":"local-user","target_host":"esxi05","target_type":"snapshot","target_id":"snapshot-456","target_name":"before-upgrade","action":"delete_snapshot","plan_id":"plan-124","result":"failed","error":"task failed: token=[REDACTED]"}
```

## Redaction rules

Audit events must not include credentials or session material. Structured fields
are rejected if they contain inline secret forms. Error text is redacted before
it is written so operators can diagnose failures without leaking passwords,
tokens, cookies, private keys, or URL userinfo.

New log files are created with private file permissions and opened for append so
existing records are preserved.
