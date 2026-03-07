## ADDED Requirements

### Requirement: FastSpec YAML export MUST be more structured than the generic draft wrapper
The FastSpec YAML export path MUST render a more structured draft artifact than the previous generic wrapper format.

#### Scenario: Structured YAML returned
- **WHEN** a reviewed draft is exported to FastSpec YAML
- **THEN** the backend writes a structured draft document
- **AND** the export remains grounded in the reviewed draft content
