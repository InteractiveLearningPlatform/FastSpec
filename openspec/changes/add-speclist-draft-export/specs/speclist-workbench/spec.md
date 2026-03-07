## ADDED Requirements

### Requirement: Workbench MUST let operators export reviewed drafts
The frontend MUST provide export controls after draft review so an operator can choose an output format and destination.

#### Scenario: Export controls shown with draft
- **WHEN** the workbench displays a reviewed draft
- **THEN** the UI shows export options
- **AND** the operator can trigger export without copying content manually

### Requirement: Workbench MUST report export outcome
The frontend MUST show the result of a completed export, including the files written by the backend.

#### Scenario: Export result displayed
- **WHEN** the backend accepts an export request
- **THEN** the frontend shows the written artifact paths
- **AND** the operator can verify where the durable files were created
