## ADDED Requirements

### Requirement: Summary command MUST support JSON output
The FastSpec CLI MUST support a machine-readable JSON mode for the `summary` command that reports the validated documents in a file or directory tree.

#### Scenario: Summary tree as JSON
- **WHEN** a contributor or agent runs `fastspec summary --json <path>`
- **THEN** the CLI emits valid JSON
- **AND** the payload includes each document's kind, identifier, title, and path

#### Scenario: Summary invalid document as JSON error
- **WHEN** `fastspec summary --json <path>` targets an invalid FastSpec document tree
- **THEN** the CLI exits with a non-zero status
- **AND** the error remains machine-consumable

### Requirement: Inspect command MUST support JSON output
The FastSpec CLI MUST support a machine-readable JSON mode for the `inspect` command that exposes parsed document metadata and document-specific details.

#### Scenario: Inspect single file as JSON
- **WHEN** a contributor or agent runs `fastspec inspect --json <file>`
- **THEN** the CLI emits valid JSON for the parsed document
- **AND** the payload includes document metadata plus kind-specific spec details

#### Scenario: Inspect tree as JSON
- **WHEN** a contributor or agent runs `fastspec inspect --json <directory>`
- **THEN** the CLI emits valid JSON for every supported document in the tree

