# Spike 0020: Pkl config guidance

GitHub issue: #20

## Question

Should Apple Pkl influence esx9s config design later, especially because the project is Go-first?

## Context

esx9s currently uses a small workstation-local YAML config for standalone ESXi hosts. The config is intentionally boring and secure: host identity, endpoint, username, external auth lookup metadata, and TLS policy. It must not store passwords, tokens, private keys, or reusable secrets.

Apple Pkl is a configuration-as-code language with validation and tooling. The Pkl project includes Go bindings and Go code generation support:

- <https://github.com/apple/pkl>
- <https://pkl-lang.org/main/current/introduction/index.html>
- <https://pkl-lang.org/go/current/quickstart.html>
- <https://pkl-lang.org/go/current/evaluation.html>
- <https://pkl-lang.org/go/current/codegen.html>

## Findings

Pkl's best ideas for esx9s are schema clarity, validation at the configuration boundary, generated typed models, and explicit control over how configuration reads external resources.

The current Pkl Go evaluator is not just a plain parser library. The Go docs describe evaluation as spawning the `pkl` CLI as a child process and communicating through message passing. That is a meaningful runtime and packaging cost for a small terminal tool whose v0.1 config should be easy to inspect, edit, test, and support offline.

Pkl code generation could become useful if esx9s config grows beyond the current host list into reusable profiles, environment overlays, defaults, or plugin-provided config modules. It could also provide a cleaner way to keep documentation, examples, and Go structs aligned. That benefit is not large enough yet to justify adding Pkl files, generated code, CLI installation requirements, or CI steps.

Security-wise, Pkl should not weaken the existing rule that esx9s config stores only references to credential sources. If Pkl is explored later, custom readers or external resource features must be constrained so config evaluation cannot casually read secrets, environment variables, local files, or network resources outside an explicit allowlist.

## Recommendation

Do not add a Pkl dependency now.

Keep YAML as the supported v0.1 user config format. Borrow Pkl's guidance instead:

- keep a typed Go config model as the source of runtime behavior
- keep validation close to load time, with clear user-facing errors
- keep examples small and schema-like
- keep secrets out of config, even when adding richer config features
- revisit generated config models only after the schema becomes large enough that hand-maintained YAML docs and Go structs begin to drift

## Later trigger points

Reopen this topic if esx9s adds:

- reusable host profiles or config inheritance
- multi-environment config overlays
- plugin-defined config sections
- generated documentation or schema artifacts
- a packaging story that can reliably provide the `pkl` CLI

Until then, Pkl should influence design discipline rather than implementation.
