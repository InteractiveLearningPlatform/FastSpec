## ADDED Requirements

### Requirement: Outline MUST identify the active draft section
The Speclist draft outline MUST identify which section is currently active during review.

#### Scenario: Reviewer navigates within a draft
- **WHEN** a reviewer has an active draft in the workbench
- **THEN** the outline shows which section is currently active

### Requirement: Outline navigation MUST update active section state
The Speclist workbench MUST update the active section when the reviewer navigates from the outline.

#### Scenario: Reviewer selects an outline entry
- **WHEN** a reviewer selects a section from the draft outline
- **THEN** that section becomes the active section
- **AND** the outline reflects the new active state

### Requirement: Draft reset MUST restore a stable active-section baseline
The Speclist workbench MUST restore active-section state to a predictable baseline when the current review session resets.

#### Scenario: Reviewer resets the draft
- **WHEN** a reviewer resets the draft to the generated baseline
- **THEN** the active section returns to the first section when one exists
