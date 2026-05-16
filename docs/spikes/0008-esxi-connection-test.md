# Spike 0008: ESXi connection test

GitHub issue: #8

## Question

What is the smallest safe connection test esx9s can offer in v0.1 for a standalone ESXi host?

## Recommendation

Build the v0.1 connection test in two layers:

1. HTTPS reachability probe against the configured vSphere SDK endpoint.
2. Authenticated govmomi login and logout using credentials supplied at runtime.

The first layer is useful as an early diagnostic because it can identify DNS/IP, TCP, TLS, proxy, and endpoint exposure problems without using a password. The second layer is the real operator-ready connection test because it proves that the endpoint accepts the selected username and password.

The connection test must remain read-only. It should not enumerate or mutate inventory during the first implementation.

## CLI shape

Initial spike stub:

```sh
ESX9S_PASSWORD='***REDACTED***' esx9s connect-test \
  --endpoint https://ESXI_HOST/sdk \
  --username operator \
  --password-env ESX9S_PASSWORD \
  --insecure-skip-verify
```

The stub performs an HTTPS GET to the SDK endpoint and reports the HTTP status. It checks whether the named password environment variable exists but does not print or use its value.

Future authenticated implementation:

```sh
ESX9S_PASSWORD='***REDACTED***' esx9s connect-test \
  --endpoint https://ESXI_HOST/sdk \
  --username operator \
  --password-env ESX9S_PASSWORD \
  --insecure-skip-verify \
  --auth-login
```

`--auth-login` should use govmomi to create a client, log in, report success/failure, and immediately log out or close the session. The command output must never include the password, session cookie, SOAP body, or Authorization headers.

## Expected output

For the reachability probe:

```text
endpoint: https://ESXI_HOST/sdk
username: operator
password: read from environment
tls_insecure_skip_verify: true
http_status: 200 OK
result: endpoint reachable; authenticated govmomi login is a follow-up implementation step
```

For the later authenticated probe:

```text
endpoint: https://ESXI_HOST/sdk
username: operator
password: read from environment
tls_insecure_skip_verify: true
sdk_reachable: true
auth_login: ok
result: connection test passed
```

## Implementation notes

- Prefer govmomi for authenticated login so esx9s follows ADR 0002.
- Use the config model from `docs/architecture/config.md`: endpoint, username, auth method, and TLS policy.
- Prompt for passwords or read them from environment variables/keychain references. Do not accept a literal password flag.
- Keep logs structured enough for diagnostics but redact credentials and session material.
- Treat `insecure_skip_verify: true` as a visible lab-only condition.
- Return non-zero exit codes for unreachable endpoints, TLS failures, missing credential sources, and failed login.

## Follow-up tasks

- Add a govmomi dependency and implement authenticated login/logout.
- Add config-file host selection, for example `connect-test --config ~/.config/esx9s/config.yaml --host lab-esxi-01`.
- Add unit tests with `httptest` for reachable, unreachable, TLS, and status reporting cases.
- Add integration-test notes for a real ESXi lab host, with credentials supplied only through environment variables or an interactive prompt.

## Spike result

The current stub proves the CLI surface and safe reachability behavior. It intentionally stops before authenticated SOAP login because that should be implemented with govmomi rather than an ad hoc HTTP/SOAP client.

Real lab probe from this workspace during the initial spike, with host details intentionally anonymized:

```sh
curl --connect-timeout 5 --max-time 10 --insecure \
  --output /dev/null \
  --write-out 'http_code=%{http_code} ssl_verify_result=%{ssl_verify_result}\n' \
  https://ESXI_HOST/sdk
```

Result: TCP connection to the lab host failed, with HTTP status `000`. No password was used, and no authenticated ESXi login was attempted.

## Issue 15 verification

Issue #15 verified the existing issue #8 spike and lightly extended the CLI instead of creating a duplicate spike.

The command now has unit coverage for:

- successful HTTPS reachability probing with `httptest`
- password environment variable presence checks without printing the value
- positive timeout validation
- explicit deferral of `--auth-login` until the govmomi implementation lands

Updated read-only lab probe:

```sh
curl --connect-timeout 5 --max-time 10 --insecure \
  --output /dev/null \
  --write-out 'http_code=%{http_code} ssl_verify_result=%{ssl_verify_result}\n' \
  https://ESXI_HOST/sdk
```

Result: TCP and TLS reached the lab host, and the server returned HTTP status `404` with certificate verification result `20` under `--insecure`. No password was used, and no authenticated ESXi login was attempted.
