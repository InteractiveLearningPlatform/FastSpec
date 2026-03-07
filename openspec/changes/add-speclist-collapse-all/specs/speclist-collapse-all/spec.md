## ADDED Requirements

### Requirement: Workbench MUST allow collapsing all draft sections
The Speclist workbench MUST allow reviewers to collapse all sections in the current draft with one action.

#### Scenario: Reviewer collapses the draft
- **WHEN** a reviewer chooses the bulk collapse action
- **THEN** every section in the current draft becomes collapsed
- **AND** the draft structure remains visible as collapsed section cards

### Requirement: Workbench MUST allow expanding all draft sections
The Speclist workbench MUST allow reviewers to expand all sections in the current draft with one action.

#### Scenario: Reviewer reopens the whole draft
- **WHEN** a reviewer chooses the bulk expand action
- **THEN** every section in the current draft becomes expanded
- **AND** section content remains unchanged

### Requirement: Bulk collapse controls MUST operate on the current draft state
Bulk collapse controls MUST reflect the currently reviewed draft rather than a stale snapshot.

#### Scenario: Draft structure changed before bulk collapse
- **WHEN** sections have been added, removed, duplicated, or reordered
- **AND** the reviewer uses a bulk collapse or expand action
- **THEN** the action applies to the current section list
