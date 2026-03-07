## ADDED Requirements

### Requirement: Outline MUST support text filtering
The Speclist draft outline MUST let reviewers filter outline entries by text.

#### Scenario: Reviewer narrows the outline by heading text
- **WHEN** a reviewer enters text into the outline filter
- **THEN** the outline shows only matching sections
- **AND** the draft content remains unchanged

### Requirement: Outline MUST support non-ready filtering
The Speclist draft outline MUST let reviewers focus on sections that still require attention.

#### Scenario: Reviewer shows only non-ready sections
- **WHEN** a reviewer enables the non-ready filter
- **THEN** the outline shows only sections with non-ready review status or review notes

### Requirement: Filtered outline entries MUST keep working navigation
Filtering the outline MUST not break section navigation.

#### Scenario: Reviewer navigates from a filtered outline
- **WHEN** a reviewer selects a filtered outline entry
- **THEN** the corresponding original draft section is focused
