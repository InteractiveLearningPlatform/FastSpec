## Why

FastSpec is positioned as an agent-first project, but the CLI currently emits only human-oriented text. That makes it harder for agents and automation to consume parsed FastSpec data without brittle output scraping.

## What Changes

- Add machine-readable JSON output for the existing `summary` and `inspect` commands.
- Define stable JSON payloads for document summaries and parsed document inspection.
- Keep the current human-readable terminal output as the default mode.
- Add tests covering JSON output and argument handling.

## Capabilities

### New Capabilities
- `json-cli-output`: Emit structured JSON from the FastSpec CLI for summary and inspect operations.

### Modified Capabilities
- `typed-spec-parsing`: Extend the typed parsing CLI surface so parsed document data can be rendered in a machine-readable format.

## Impact

- Affected code: `apps/fastspec-cli`, `crates/fastspec-core`, `crates/fastspec-model`
- Affected docs: CLI usage docs
- Dependencies: Rust JSON serialization support via `serde`
