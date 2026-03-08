## ADDED Requirements

### Requirement: Workbench MUST support sequential outline navigation
The Speclist workbench MUST let reviewers move to the next or previous section from the draft outline.

#### Scenario: Reviewer steps forward through the draft
- **WHEN** a reviewer chooses the next navigation action
- **THEN** the next section in the current outline view becomes active

### Requirement: Sequential navigation MUST respect outline filters
The Speclist workbench MUST base sequential navigation on the currently visible outline entries.

#### Scenario: Reviewer navigates within a filtered outline
- **WHEN** outline filters are active
- **AND** the reviewer uses next or previous navigation
- **THEN** navigation follows only the currently visible outline entries

### Requirement: Sequential navigation MUST keep working when the active section is filtered out
The Speclist workbench MUST recover predictably when the active section is not visible in the current outline view.

#### Scenario: Active section is not in the filtered outline
- **WHEN** the reviewer uses sequential navigation
- **AND** the active section is not part of the filtered outline
- **THEN** navigation selects a visible outline entry and focuses it
