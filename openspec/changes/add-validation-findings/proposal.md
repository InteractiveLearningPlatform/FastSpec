## Why

FastSpec can now parse documents and emit machine-readable summaries, but it still has no way to report cross-document validation findings beyond parse failure. Agents need explicit findings they can inspect, triage, and fix before any generation-oriented workflow becomes useful.

## What Changes

- Add a `validate` command to the FastSpec CLI.
- Introduce structured validation findings for duplicate identifiers and invalid references across a spec tree.
- Support both human-readable and JSON output for validation results.
- Add tests covering successful validation and failing trees with actionable findings.

## Capabilities

### New Capabilities
- `validation-findings`: Report structured validation findings for FastSpec trees instead of only parse success or failure.

### Modified Capabilities
- `typed-spec-parsing`: Extend the CLI parsing surface with a validation command that reuses the typed model.
- `json-cli-output`: Expose validation findings in machine-readable JSON form.

## Impact

- Affected code: `apps/fastspec-cli`, `crates/fastspec-core`, `crates/fastspec-model`
- Affected examples/tests: validation fixtures and CLI integration tests
- No new external dependencies are required beyond the current serde-based stack
