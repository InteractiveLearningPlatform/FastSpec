## ADDED Requirements

### Requirement: Workbench MUST allow section collapse
The Speclist workbench MUST allow reviewers to collapse a draft section during review.

#### Scenario: Reviewer collapses a section
- **WHEN** a reviewer collapses a draft section
- **THEN** the section remains visible in the draft list
- **AND** only compact section information remains visible

### Requirement: Workbench MUST allow collapsed sections to expand again
The Speclist workbench MUST allow a collapsed draft section to expand back into editable form.

#### Scenario: Reviewer reopens a collapsed section
- **WHEN** a reviewer expands a collapsed draft section
- **THEN** the section editing fields become visible again
- **AND** the section content remains unchanged

### Requirement: Collapse state MUST track local section review changes
The workbench MUST keep collapse state aligned with the current in-memory section list during review.

#### Scenario: Section order or count changes
- **WHEN** a reviewer reorders, duplicates, removes, or resets sections
- **THEN** collapse state remains aligned with the resulting section list
- **AND** reset restores a fully expanded baseline
