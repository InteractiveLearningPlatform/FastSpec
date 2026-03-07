## ADDED Requirements

### Requirement: Workbench MUST show a draft outline
The Speclist workbench MUST show a compact outline for the current draft during review.

#### Scenario: Reviewer inspects the draft structure
- **WHEN** a draft is present in the workbench
- **THEN** the reviewer can see an ordered list of the draft sections
- **AND** each outline entry identifies the section it represents

### Requirement: Outline MUST support section navigation
The draft outline MUST let reviewers navigate to a specific section in the current draft.

#### Scenario: Reviewer jumps to a section
- **WHEN** a reviewer selects a section from the outline
- **THEN** the corresponding draft section is brought into view

### Requirement: Outline navigation MUST reopen collapsed target sections
The draft outline MUST reopen a collapsed section when the reviewer navigates to it.

#### Scenario: Reviewer opens a collapsed section from the outline
- **WHEN** the target section is collapsed
- **AND** the reviewer selects it from the outline
- **THEN** that section becomes expanded
- **AND** the rest of the draft remains unchanged
