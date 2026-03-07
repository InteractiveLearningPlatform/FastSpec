## ADDED Requirements

### Requirement: Speclist MUST export to active OpenSpec change artifacts
The system MUST support exporting a reviewed draft directly into supported artifact paths of an active OpenSpec change.

#### Scenario: Export to proposal artifact
- **WHEN** an operator exports a reviewed draft to an active OpenSpec change proposal target
- **THEN** the backend writes the draft to that change's `proposal.md`
- **AND** it writes a citation sidecar next to the exported markdown

#### Scenario: Export to capability spec artifact
- **WHEN** an operator exports a reviewed draft to an active OpenSpec change spec target with a capability name
- **THEN** the backend writes the draft to `specs/<capability>/spec.md`
- **AND** it returns the written artifact paths

### Requirement: System MUST list active OpenSpec changes for export targeting
The backend MUST provide the workbench with the set of active OpenSpec changes that can accept repo-aware export.

#### Scenario: Active changes returned
- **WHEN** the workbench requests available OpenSpec change targets
- **THEN** the backend returns active changes only
- **AND** each change includes the supported artifact target types
