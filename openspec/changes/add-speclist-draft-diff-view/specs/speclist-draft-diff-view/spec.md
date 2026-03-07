## ADDED Requirements

### Requirement: Workbench MUST preserve the original generated draft for comparison
The Speclist workbench MUST keep the originally generated draft available while reviewers edit the live draft.

#### Scenario: Reviewer edits generated draft
- **WHEN** a reviewer generates a draft and then modifies it
- **THEN** the original generated draft remains available in the workbench
- **AND** the current draft can be compared against that original state

### Requirement: Workbench MUST show structured draft differences
The Speclist workbench MUST show the differences between the current edited draft and the original generated draft before export.

#### Scenario: Reviewer inspects draft diff
- **WHEN** the current draft differs from the original generated draft
- **THEN** the workbench shows changed title, summary, and section content
- **AND** added or removed sections are clearly identified
