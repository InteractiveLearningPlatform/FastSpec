---
area: repository
component: bootstrap
kind: proposal
status: active
tags:
  - openspec
  - fastspec
  - bootstrap
---

# Bootstrap Agent-First Repo

## Problem

FastSpec currently exists as a short manifesto only. There is no operational OpenSpec workflow, no durable project structure, no reusable YAML templates, and no example app specification that shows how FastSpec should be used in practice.

That makes the project hard for agents to contribute to and hard for humans to evaluate. Every future implementation step would otherwise start by re-inventing repo layout, artifact conventions, and example inputs.

## Proposed Change

Bootstrap the repository as an agent-first project with:

- OpenSpec initialized and documented as the default change workflow
- repo documentation that explains the FastSpec/OpenSpec relationship
- reusable YAML templates for project, module, and agent capability specs
- a first example app under `examples/archlint-reproduction/`
- a minimal Rust workspace under `apps/` and `crates/` that can inspect example specs

## Impact

- Affected areas: OpenSpec workflow, docs, templates, examples, Rust runtime scaffolding
- No generation engine or MCP runtime is implemented yet; this is a foundation slice for later Rust work
- The repo becomes immediately usable for spec authoring and future implementation slices

## Non-Goals

- Implementing the full Rust runtime
- Building actual MCP servers
- Generating runnable applications from FastSpec documents in this slice
