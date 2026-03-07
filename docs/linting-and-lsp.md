---
area: docs
component: developer-experience
kind: reference
status: active
language: rust
framework: lsp
tags:
  - rustfmt
  - clippy
  - rust-analyzer
  - lsp
---

# Linting And LSP Context

FastSpec currently targets Rust for its runtime and core libraries, so the default developer-experience stack should be Rust-native.

## Language Server Protocol

The Language Server Protocol defines a common protocol between editors and language servers. In practice, it is the abstraction that makes semantic navigation, completion, diagnostics, and refactoring portable across editors.

For FastSpec, LSP matters for two reasons:

- human contributors need strong code intelligence in whatever editor they use
- agent tooling benefits from semantic rather than purely text-based navigation

Treat LSP support as a first-class capability, not an optional editor convenience.

## rust-analyzer

rust-analyzer is the default Rust language-server reference.

- It provides go-to-definition, find-all-references, refactorings, and code completion.
- It integrates formatting with rustfmt and diagnostics with rustc and Clippy.
- It works in any editor that supports LSP.

For FastSpec, rust-analyzer should be the baseline assumption for local editor support and for any future agent tooling that can consume LSP-backed data.

## rustfmt

rustfmt is the default formatter.

- Its standard workspace entrypoint is `cargo fmt`.
- It supports `--check` for CI enforcement.
- Project behavior can be pinned with `rustfmt.toml`.

FastSpec should prefer stable, repo-checked formatting rules over editor-local formatting drift.

## Clippy

Clippy is the default linter.

- Rust documents Clippy as a collection of lints for catching common mistakes and improving code quality.
- It includes categories for correctness, suspicious code, style, complexity, and performance, with stricter groups such as `pedantic` available when needed.

For FastSpec, Clippy should be the default semantic lint gate before introducing third-party Rust lint layers.

## Practical Direction For FastSpec

- Use rustfmt for formatting.
- Use Clippy for linting.
- Use rust-analyzer as the editor and agent semantic backend.
- Keep LSP-friendly project structure and avoid build setups that make analysis brittle.

If FastSpec later adds Python, YAML, or Markdown-heavy tooling, then complementary linters can be layered in, but the Rust path should stay the default for core runtime work.
