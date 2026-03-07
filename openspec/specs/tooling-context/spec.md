# tooling-context Specification

## Purpose
TBD - created by archiving change add-tooling-context. Update Purpose after archive.
## Requirements
### Requirement: Repo Maintains Durable Tooling Context

The repository MUST keep durable context for the core workflow and tooling stack it expects contributors and agents to use.

#### Scenario: Contributor needs repo-local tooling guidance

- **WHEN** a contributor or agent needs context about the workflow or supporting tools
- **THEN** the repository provides focused docs for the relevant tooling
- **AND** the context is stored in repo-local files rather than only in chat history

### Requirement: OpenSpec Configuration Includes Compact Tooling Summary

The repository MUST keep a compact summary of its external tooling assumptions in `openspec/config.yaml`.

#### Scenario: Agent creates a new change

- **WHEN** an agent reads project context from `openspec/config.yaml`
- **THEN** it can see the preferred workflow and external tooling assumptions
- **AND** it does not need to load every detailed doc by default

### Requirement: Inspiration Notes Are Kept Separate From Current Behavior

The repository MUST separate inspiration from current implementation commitments.

#### Scenario: Repo borrows ideas from another project

- **WHEN** the repository documents ideas from another tool or framework
- **THEN** it stores them as inspiration or reference notes
- **AND** it avoids implying that those behaviors already exist in FastSpec

