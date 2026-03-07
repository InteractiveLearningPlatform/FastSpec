## ADDED Requirements

### Requirement: Draft generation MUST support preset review structures
The system MUST support a small set of preset draft structures so different review modes can start from more intentional grounded sections.

#### Scenario: Proposal preset generated
- **WHEN** a reviewer generates a draft with the `proposal` preset
- **THEN** the resulting draft uses proposal-oriented sections
- **AND** the draft payload records that the preset was `proposal`

#### Scenario: Requirements preset generated
- **WHEN** a reviewer generates a draft with the `requirements` preset
- **THEN** the resulting draft includes requirements-oriented sections
- **AND** the generated draft remains editable before export

### Requirement: Workbench MUST expose draft preset selection
The Speclist workbench MUST let reviewers choose a draft preset before generation.

#### Scenario: Reviewer chooses preset from the UI
- **WHEN** a reviewer selects a preset in the workbench and generates a draft
- **THEN** the frontend sends the selected preset with the draft request
- **AND** the loaded draft shows which preset was used
