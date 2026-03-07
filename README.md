# FastSpec

FastSpec is a YAML-first, agent-oriented layer for spec-driven development.
It pairs well with OpenSpec:

- OpenSpec manages change workflow: proposal, specs, design, tasks, archive.
- FastSpec defines compact, reusable domain specs that agents can retrieve and compose.

The target audience is AI agents first and humans second. Human readability still matters, but the default bias is toward structured, low-token artifacts that travel well across tools.

## Principles

- Made for agents, with humane interfaces layered on top
- Fast and efficient before ornate
- Decentralized by default
- As complex as needed to solve the task; lightweight is a preference, not a dogma
- Formal structure where it improves machine reasoning

## Direction

- Primary implementation language: Rust
- Primary spec format: YAML
- Supporting docs: Markdown with Mermaid diagrams when diagrams help
- ML and research tooling: Python, with room for JAX-based experiments

Platform modules:

- Memory: durable project and agent memory
- Spec designer: compact, context-aware spec authoring
- Swarm manager: multi-agent task splitting and coordination
- Knowledge base: retrieval-oriented technical memory
- Speclist: ingestion and RAG workbench for converting existing documentation into grounded specs

## Repository Layout

- `openspec/` - change workflow and active implementation artifacts
- `docs/` - human-readable architecture, workflow, and tooling-context docs
- `templates/` - reusable FastSpec YAML templates
- `examples/` - end-to-end example specs, including app creation examples
- `apps/` - executable surfaces such as the future CLI
- `apps/speclist-api/` - Go microservice for ingestion, retrieval, and grounded draft generation
- `apps/speclist-web/` - React workbench for operators using Speclist
- `crates/` - reusable Rust libraries

Useful starting docs:

- `docs/working-with-openspec.md`
- `docs/tooling-stack.md`
- `docs/linting-and-lsp.md`
- `docs/opencode-ideas.md`
- `docs/speclist.md`

## Getting Started

1. Start from OpenSpec:
   `openspec list --json`
2. Create or continue a change:
   `/opsx:propose "idea"` then `/opsx:apply`
3. Use the templates in `templates/` to draft durable YAML specs.
4. Use `examples/archlint-reproduction/` as the initial app-creation reference.

## Pre-Commit

This repo uses `pre-commit` for Rust checks before commit.

1. Install `pre-commit`.
2. Run `pre-commit install`.
3. Run `pre-commit run --all-files` to validate the current tree.

The default Rust hooks are:

- `cargo fmt --all -- --check`
- `cargo clippy --all-targets --all-features -- -D warnings`

Rust formatting is configured in `rustfmt.toml` with a `140` character line width and repo-wide defaults for Rust 2024.

## CI/CD

GitHub Actions now validates each stack independently on pushes and pull requests to `main` and `develop`:

- Rust: format, clippy, and workspace tests
- Go: `go test ./...` for `apps/speclist-api`
- Web: `npm ci` and `npm run build` for `apps/speclist-web`

To publish release artifacts, push a version tag such as `v0.1.0`. The release workflow packages:

- the `fastspec` Linux CLI binary
- the `speclist-api` Linux service binary
- the built `speclist-web` static bundle

## Current Scope

This repository now includes:

- the OpenSpec workflow and archived change history
- the Rust FastSpec model, validation, graph, plan, and scaffold runtime
- the first `speclist` service slice for DOCX/Confluence ingestion, spec retrieval, and grounded draft review
