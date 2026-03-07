## Why

FastSpec currently treats YAML documents as raw text and only detects the `kind:` line. That is enough for bootstrap validation, but not enough for meaningful inspection, validation, or future generation work.

## What Changes

- Add typed Rust models for the current FastSpec document types.
- Add YAML parsing that loads FastSpec files into typed document structures instead of string-only kind detection.
- Extend validation so the CLI can report document metadata and fail on malformed required fields.
- Expand the CLI beyond `summary` so contributors can inspect parsed FastSpec documents directly.
- Add tests covering successful parsing and invalid-document failures.

## Capabilities

### New Capabilities
- `typed-spec-parsing`: Parse FastSpec YAML documents into typed Rust structures and expose them through the CLI.

### Modified Capabilities
- `repository-foundation`: The minimal executable Rust scaffold will evolve from kind-only inspection to typed parsing and validation of example specs.

## Impact

- Affected code: `apps/fastspec-cli`, `crates/fastspec-core`, `crates/fastspec-model`
- Affected examples: `examples/archlint-reproduction/specs`
- Dependencies: likely add a YAML parsing crate such as `serde`/`serde_yaml`
