## ADDED Requirements

### Requirement: Workbench MUST support keyboard shortcuts for sequential outline navigation
The Speclist workbench MUST let reviewers trigger sequential outline navigation from the keyboard.

#### Scenario: Reviewer moves to the next section from the keyboard
- **WHEN** a reviewer uses the next-section keyboard shortcut during draft review
- **THEN** the next visible outline section becomes active

### Requirement: Keyboard shortcuts MUST not interfere with text editing
The Speclist workbench MUST ignore navigation shortcuts while the reviewer is typing into form controls.

#### Scenario: Reviewer types inside the draft editor
- **WHEN** keyboard focus is inside an input, textarea, or select
- **THEN** outline navigation shortcuts do not trigger

### Requirement: Keyboard shortcuts MUST follow the current outline view
The Speclist workbench MUST apply keyboard navigation to the currently visible outline entries.

#### Scenario: Reviewer uses shortcuts with active outline filters
- **WHEN** outline filters are active
- **AND** the reviewer uses an outline navigation shortcut
- **THEN** navigation follows the current filtered outline entries
