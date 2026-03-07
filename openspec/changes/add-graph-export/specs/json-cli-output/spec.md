## MODIFIED Requirements

### Requirement: Summary command MUST support JSON output
The FastSpec CLI MUST support machine-readable JSON output for `summary`, `validate`, and `graph` commands so automation can consume overview, validation, and graph data directly.

#### Scenario: Summary tree as JSON
- **WHEN** a contributor or agent runs `fastspec summary --json <path>`
- **THEN** the CLI emits valid JSON
- **AND** the payload includes each document's kind, identifier, title, and path

#### Scenario: Summary invalid document as JSON error
- **WHEN** `fastspec summary --json <path>` targets an invalid FastSpec document tree
- **THEN** the CLI exits with a non-zero status
- **AND** the error remains machine-consumable

#### Scenario: Validation findings as JSON
- **WHEN** a contributor or agent runs `fastspec validate --json <path>`
- **THEN** the CLI emits valid JSON
- **AND** the payload includes a validity result plus any findings

#### Scenario: Graph export as JSON
- **WHEN** a contributor or agent runs `fastspec graph --json <path>`
- **THEN** the CLI emits valid JSON
- **AND** the payload includes normalized graph nodes and edges
