## ADDED Requirements

### Requirement: Workbench MUST allow draft section reordering
The Speclist workbench MUST allow reviewers to change draft section order during review.

#### Scenario: Reviewer moves a section up
- **WHEN** a reviewer moves a section upward in the draft editor
- **THEN** the section appears earlier in the draft
- **AND** the section content remains unchanged

#### Scenario: Reviewer moves a section down
- **WHEN** a reviewer moves a section downward in the draft editor
- **THEN** the section appears later in the draft
- **AND** the section content remains unchanged

### Requirement: Review metadata MUST stay aligned after reorder
The workbench MUST keep section review metadata aligned when sections are reordered.

#### Scenario: Reorder preserves section review state
- **WHEN** a reviewer reorders sections that already have review flags
- **THEN** the flags move with the corresponding sections
- **AND** the review summary still refers to the correct section content
