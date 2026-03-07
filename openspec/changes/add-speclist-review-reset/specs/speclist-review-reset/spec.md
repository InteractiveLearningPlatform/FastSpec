## ADDED Requirements

### Requirement: Workbench MUST allow review reset to the generated baseline
The Speclist workbench MUST allow reviewers to reset the current draft back to the original generated draft.

#### Scenario: Reviewer resets edited draft
- **WHEN** a reviewer resets the current draft
- **THEN** the editable draft matches the original generated draft again
- **AND** the original generated snapshot remains available for future comparison

### Requirement: Reset MUST clear review-only state
The workbench MUST clear review-only state that no longer applies after reset.

#### Scenario: Reset clears flags and export result
- **WHEN** a reviewer resets the current draft
- **THEN** section review flags are cleared
- **AND** previous export result state is cleared
