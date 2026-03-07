## ADDED Requirements

### Requirement: Workbench MUST allow section duplication
The Speclist workbench MUST allow reviewers to duplicate a draft section during review.

#### Scenario: Reviewer duplicates a section
- **WHEN** a reviewer duplicates a draft section
- **THEN** a copy of the section is inserted into the draft
- **AND** the copied section content matches the source section

### Requirement: Duplicated sections MUST preserve review metadata
The workbench MUST preserve section review metadata when duplicating a section.

#### Scenario: Duplicate keeps review flag state
- **WHEN** a reviewer duplicates a section that has review status or notes
- **THEN** the duplicated section starts with the same review metadata
- **AND** both sections remain independently editable afterward
