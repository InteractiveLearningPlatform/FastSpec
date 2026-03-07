## 1. JSON serialization support

- [x] 1.1 Add JSON serialization support for the CLI-facing summary and inspect data structures
- [x] 1.2 Define stable JSON payload shapes for summary output and parsed document inspection

## 2. CLI output mode

- [x] 2.1 Extend `fastspec summary` and `fastspec inspect` to accept a `--json` flag
- [x] 2.2 Keep the current human-readable output as the default while returning machine-readable JSON when requested

## 3. Verification and docs

- [x] 3.1 Add tests covering JSON output and invalid-input behavior
- [x] 3.2 Update CLI docs to describe the new JSON mode and example usage
