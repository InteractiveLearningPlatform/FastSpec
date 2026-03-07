## MODIFIED Requirements

### Requirement: Agent-First Project Layout

The repository MUST provide a layout that makes OpenSpec changes, durable FastSpec YAML artifacts, templates, examples, and future runtime code easy to discover, including service-oriented products that use a different implementation stack when the product boundary justifies it.

#### Scenario: Bootstrap layout exists

- **WHEN** a contributor opens the repository root
- **THEN** they find `openspec/`, `docs/`, `templates/`, `examples/`, `apps/`, and `crates/`
- **AND** each area explains or demonstrates its intended role

#### Scenario: Speclist service layout exists

- **WHEN** the repository adds the `speclist` product
- **THEN** contributors can discover separate backend and frontend application roots
- **AND** the layout makes it clear that the Go backend and React frontend belong to the same product surface
