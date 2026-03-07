## ADDED Requirements

### Requirement: CLI MUST export an ordered implementation plan
The FastSpec CLI MUST provide a `plan` command that exports an ordered implementation plan derived from a validation-clean FastSpec graph.

#### Scenario: Export plan from valid tree
- **WHEN** a contributor or agent runs `fastspec plan <path>` on a valid FastSpec tree
- **THEN** the CLI exits successfully
- **AND** it emits a deterministic ordered plan

#### Scenario: Reject plan export from invalid tree
- **WHEN** a contributor or agent runs `fastspec plan <path>` on a FastSpec tree with validation findings
- **THEN** the CLI exits with a non-zero status
- **AND** it reports that plan export requires a validation-clean tree

### Requirement: Plan MUST include explicit step dependencies
The exported plan MUST represent steps with explicit identifiers, phases, and dependency references.

#### Scenario: Module steps depend on prerequisites
- **WHEN** one module depends on another module in the same project
- **THEN** the dependent module step references the prerequisite module step as a dependency

#### Scenario: Workflow steps follow structural steps
- **WHEN** the project defines workflows
- **THEN** the exported plan includes workflow-oriented steps after the project and module structure steps

### Requirement: Plan MUST support JSON output
The `plan` command MUST support machine-readable JSON output.

#### Scenario: Plan as JSON
- **WHEN** a contributor or agent runs `fastspec plan --json <path>`
- **THEN** the CLI emits valid JSON containing the ordered plan steps
