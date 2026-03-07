## ADDED Requirements

### Requirement: Backend MUST resolve source ids to full source detail
The system MUST resolve a source id to the matching stored source document so operators can inspect source metadata and chunks directly.

#### Scenario: Source detail lookup succeeds
- **WHEN** an operator requests detail for a source id that exists in the corpus
- **THEN** the backend returns the full source document
- **AND** the response includes the source metadata and chunk list

#### Scenario: Source detail lookup fails
- **WHEN** an operator requests detail for a source id that does not exist
- **THEN** the backend returns an error
- **AND** no partial source payload is returned

### Requirement: Workbench MUST expose source detail review
The Speclist workbench MUST allow operators to inspect indexed sources directly from the source list.

#### Scenario: Operator opens a source detail panel
- **WHEN** an operator chooses a source from the indexed source list
- **THEN** the workbench loads the source detail
- **AND** shows the source title, location, metadata, and chunk inventory in a dedicated panel
