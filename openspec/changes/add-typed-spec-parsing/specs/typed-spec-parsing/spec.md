## ADDED Requirements

### Requirement: FastSpec documents MUST parse into typed runtime models
The runtime MUST parse supported FastSpec YAML documents into typed Rust structures for project, module, and agent capability specs.

#### Scenario: Parse project spec
- **WHEN** the runtime loads a valid `ProjectSpec` YAML document
- **THEN** it returns a typed project document structure with metadata and project-specific fields populated

#### Scenario: Parse module spec
- **WHEN** the runtime loads a valid `ModuleSpec` YAML document
- **THEN** it returns a typed module document structure with metadata and module-specific fields populated

#### Scenario: Parse agent capability spec
- **WHEN** the runtime loads a valid `AgentCapabilitySpec` YAML document
- **THEN** it returns a typed agent capability document structure with metadata and capability-specific fields populated

### Requirement: Validation MUST reject malformed supported documents
The runtime MUST reject supported FastSpec YAML documents when required fields are missing or when the declared kind does not match the document structure being parsed.

#### Scenario: Missing required field
- **WHEN** a supported FastSpec YAML document omits a required field
- **THEN** validation fails with an error that identifies the invalid document

#### Scenario: Unknown kind
- **WHEN** a FastSpec YAML document declares an unsupported `kind`
- **THEN** validation fails instead of silently treating the document as valid

### Requirement: CLI MUST expose parsed document details
The CLI MUST provide a command that renders parsed FastSpec document details from a file or tree so contributors can inspect typed data without reading implementation code.

#### Scenario: Inspect a single document
- **WHEN** a contributor runs the CLI inspection command on a valid FastSpec YAML file
- **THEN** the CLI prints the document kind, identifier, and core metadata derived from the typed model

#### Scenario: Inspect a spec tree
- **WHEN** a contributor runs the CLI inspection command on a directory of FastSpec YAML files
- **THEN** the CLI prints parsed details for each supported document in the tree
