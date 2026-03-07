## ADDED Requirements

### Requirement: Workbench MUST allow section-level review flags
The Speclist workbench MUST allow reviewers to assign a lightweight review status to each draft section.

#### Scenario: Reviewer flags a section
- **WHEN** a reviewer marks a draft section during review
- **THEN** the workbench stores the selected review status for that section
- **AND** the section continues to be editable

### Requirement: Workbench MUST summarize flagged sections before export
The Speclist workbench MUST show a compact summary of section flags before export.

#### Scenario: Reviewer checks unresolved sections
- **WHEN** one or more sections are marked as `needs-work` or `blocked`
- **THEN** the workbench shows those flagged sections in a review summary
- **AND** optional reviewer notes remain visible with the flagged section
