---
area: repository
component: foundation
kind: spec
status: proposed
tags:
  - openspec
  - fastspec
  - yaml
---

## ADDED Requirements

### Requirement: Agent-First Project Layout

The repository MUST provide a layout that makes OpenSpec changes, durable FastSpec YAML artifacts, templates, examples, and future runtime code easy to discover.

#### Scenario: Bootstrap layout exists

- **WHEN** a contributor opens the repository root
- **THEN** they find `openspec/`, `docs/`, `templates/`, `examples/`, `apps/`, and `crates/`
- **AND** each area explains or demonstrates its intended role

### Requirement: OpenSpec Is The Default Change Workflow

The repository MUST document OpenSpec as the default workflow for proposing, implementing, and archiving work.

#### Scenario: Contributor needs to start work

- **WHEN** a contributor reads the top-level docs
- **THEN** they can discover the OpenSpec flow
- **AND** they can identify where active change artifacts live

### Requirement: Durable YAML Templates Exist

The repository MUST provide reusable YAML templates for core FastSpec document types.

#### Scenario: Author creates a new system spec

- **WHEN** an author needs to describe a project or module
- **THEN** the repository provides template YAML files
- **AND** the templates reflect a compact, retrieval-friendly structure

### Requirement: Example Application Specs Exist

The repository MUST include at least one realistic example that demonstrates how FastSpec documents can describe an application.

#### Scenario: Agent needs a concrete reference

- **WHEN** an agent or contributor looks for an example
- **THEN** they can use `examples/archlint-reproduction/`
- **AND** the example includes a project-level spec plus module-level specs

### Requirement: Minimal Executable Rust Scaffold Exists

The repository MUST include a minimal Rust workspace that can inspect FastSpec example documents.

#### Scenario: Contributor validates example specs

- **WHEN** a contributor runs the workspace tests or CLI
- **THEN** the workspace can detect the example spec kinds
- **AND** it can summarize the example tree without external services
