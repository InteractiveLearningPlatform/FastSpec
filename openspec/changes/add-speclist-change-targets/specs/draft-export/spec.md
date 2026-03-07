## ADDED Requirements

### Requirement: Export MUST support OpenSpec change targets
The export workflow MUST support repo-aware OpenSpec change targets in addition to generic filesystem destinations.

#### Scenario: OpenSpec target selected
- **WHEN** the operator selects an OpenSpec change target for export
- **THEN** the backend resolves the artifact path from the change name and artifact type
- **AND** the export response reports the resulting written files
