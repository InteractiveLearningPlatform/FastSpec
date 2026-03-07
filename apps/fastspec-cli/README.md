# fastspec-cli

Minimal Rust CLI entrypoint.

Current commands:

- `fastspec summary <path>` to validate and summarize a FastSpec tree
- `fastspec inspect <path>` to print typed document details for one FastSpec file or a directory tree
- `fastspec validate <path>` to report structured validation findings for a FastSpec file or tree
- `fastspec graph <path>` to export a normalized project graph from a validation-clean FastSpec tree
- `fastspec plan <path>` to export an ordered implementation plan from a validation-clean FastSpec graph
- `fastspec summary --json <path>` to emit machine-readable document summaries
- `fastspec inspect --json <path>` to emit machine-readable parsed document details
- `fastspec validate --json <path>` to emit machine-readable validation findings
- `fastspec graph --json <path>` to emit machine-readable graph export output
- `fastspec plan --json <path>` to emit machine-readable plan export output

Planned next:

- generation-oriented commands
