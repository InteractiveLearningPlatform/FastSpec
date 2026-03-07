## ADDED Requirements

### Requirement: Workbench MUST summarize export readiness
The Speclist workbench MUST summarize whether the current draft is ready for export.

#### Scenario: Draft has blockers
- **WHEN** the current draft has blocking review status or missing required fields
- **THEN** the workbench shows a blocking export readiness state
- **AND** lists the blocking reasons before export

#### Scenario: Draft has warnings only
- **WHEN** the current draft has warnings but no blockers
- **THEN** the workbench shows a warning readiness state
- **AND** the reviewer can still inspect the warning list before export

### Requirement: Workbench MUST keep readiness near export controls
The Speclist workbench MUST present export readiness where reviewers decide whether to export.

#### Scenario: Reviewer prepares export
- **WHEN** a draft is visible in the workbench
- **THEN** the export area includes a readiness summary
- **AND** blockers and warnings are visible without leaving the draft review flow
