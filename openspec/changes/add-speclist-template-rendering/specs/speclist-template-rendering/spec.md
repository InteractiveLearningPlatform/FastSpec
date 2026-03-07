## ADDED Requirements

### Requirement: OpenSpec-targeted export MUST render artifact-aware templates
The system MUST render reviewed drafts into different content templates based on the selected OpenSpec artifact target.

#### Scenario: Proposal target uses proposal template
- **WHEN** an operator exports a reviewed draft to an OpenSpec proposal target
- **THEN** the rendered markdown follows a proposal-oriented structure
- **AND** it is not emitted as a generic draft body dump

#### Scenario: Tasks target uses checklist template
- **WHEN** an operator exports a reviewed draft to an OpenSpec tasks target
- **THEN** the rendered markdown uses task checklist formatting
- **AND** the result is shaped for later refinement into implementation tasks

### Requirement: Spec target MUST render requirement-oriented output
The system MUST render OpenSpec spec targets into requirement-oriented markdown with scenario structure.

#### Scenario: Capability spec exported
- **WHEN** an operator exports a reviewed draft to `specs/<capability>/spec.md`
- **THEN** the rendered markdown contains requirement-oriented sections
- **AND** it remains suitable for later manual refinement into normative spec language
