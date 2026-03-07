# Repository Structure

- `openspec/`
  OpenSpec workflow, active changes, and future archived specs.
- `docs/`
  Concise explanations for humans and retrieval entry points for agents, including durable tooling context.
- `templates/`
  Reusable YAML starters for FastSpec documents.
- `examples/`
  Concrete end-to-end references, including app-creation examples.
- `apps/`
  Executable surfaces, including the Rust FastSpec CLI and the Speclist product apps.
- `crates/`
  Reusable Rust libraries for the FastSpec runtime.

Current mixed-stack product exception:

- `apps/speclist-api/`
  Go hexagonal microservice for document ingestion, spec indexing, retrieval, and draft generation.
- `apps/speclist-web/`
  React workbench for importing sources, searching grounded context, and reviewing drafts.

The repo is intentionally structured so product-specific services can live next to the Rust FastSpec runtime without blurring their boundaries.
