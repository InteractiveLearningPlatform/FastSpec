## ADDED Requirements

### Requirement: Speclist MUST export reviewed drafts into durable files
The system MUST export a reviewed draft into a caller-selected target directory on disk instead of limiting the draft to in-memory or UI-only review.

#### Scenario: Draft exported to filesystem
- **WHEN** an operator exports a reviewed draft
- **THEN** the backend writes one or more files to the requested target directory
- **AND** the export response reports the written artifact paths

### Requirement: Export MUST support OpenSpec markdown and FastSpec YAML
The system MUST support at least one OpenSpec-oriented markdown export format and one FastSpec-oriented YAML export format.

#### Scenario: Export OpenSpec markdown
- **WHEN** the operator selects the OpenSpec markdown export format
- **THEN** the backend writes a markdown artifact suitable for change-time workflow usage

#### Scenario: Export FastSpec YAML
- **WHEN** the operator selects the FastSpec YAML export format
- **THEN** the backend writes a YAML artifact suitable for durable spec refinement

### Requirement: Export MUST preserve machine-readable citations
The system MUST write machine-readable citation metadata alongside the exported draft so agents can recover provenance without reparsing prose.

#### Scenario: Citation sidecar written
- **WHEN** the backend exports a draft
- **THEN** it writes a machine-readable citation sidecar file next to the primary artifact
- **AND** the sidecar identifies the draft sections and supporting citations

### Requirement: Export MUST avoid silent overwrite
The system MUST reject export requests that would overwrite existing export files for the selected target name.

#### Scenario: Existing export blocks write
- **WHEN** an export request targets an existing artifact path
- **THEN** the backend returns an error
- **AND** it does not partially overwrite the existing export set
