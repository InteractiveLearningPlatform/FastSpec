## ADDED Requirements

### Requirement: Export MUST preserve useful structure when targeting OpenSpec artifacts
The export workflow MUST preserve useful structure by rendering reviewed drafts into artifact-aware templates for supported OpenSpec targets.

#### Scenario: Structured export returned
- **WHEN** a reviewed draft is exported to a supported OpenSpec target
- **THEN** the backend writes a structured artifact template
- **AND** the export remains grounded in the reviewed draft content
