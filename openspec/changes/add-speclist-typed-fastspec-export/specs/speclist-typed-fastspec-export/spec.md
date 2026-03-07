## ADDED Requirements

### Requirement: FastSpec YAML export MUST render a typed draft document
The system MUST render FastSpec YAML export as a typed draft document rather than a single generic blob-style export wrapper.

#### Scenario: Typed YAML draft exported
- **WHEN** an operator exports a reviewed draft using the FastSpec YAML format
- **THEN** the backend writes a typed YAML draft document
- **AND** the document preserves structured sections for later refinement

### Requirement: FastSpec YAML export MUST preserve structured requirement content
The typed YAML draft export MUST preserve requirement-oriented content from the reviewed draft in structured YAML fields.

#### Scenario: Requirement content represented in YAML
- **WHEN** the reviewed draft contains requirement-oriented content
- **THEN** the exported YAML stores that content in stable structured fields
- **AND** the result remains suitable for durable spec refinement
