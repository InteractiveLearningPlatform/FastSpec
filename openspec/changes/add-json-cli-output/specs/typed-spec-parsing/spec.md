## MODIFIED Requirements

### Requirement: CLI MUST expose parsed document details
The CLI MUST provide a command that renders parsed FastSpec document details from a file or tree so contributors can inspect typed data in either human-readable or machine-readable form.

#### Scenario: Inspect a single document
- **WHEN** a contributor runs the CLI inspection command on a valid FastSpec YAML file
- **THEN** the CLI prints or emits the document kind, identifier, and core metadata derived from the typed model

#### Scenario: Inspect a spec tree
- **WHEN** a contributor runs the CLI inspection command on a directory of FastSpec YAML files
- **THEN** the CLI prints or emits parsed details for each supported document in the tree
