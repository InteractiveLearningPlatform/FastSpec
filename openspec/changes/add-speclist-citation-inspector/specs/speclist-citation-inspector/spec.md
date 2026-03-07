## ADDED Requirements

### Requirement: Backend MUST resolve citation strings to grounded source context
The system MUST resolve a citation string to the matching source document and chunk so reviewers can inspect the grounding behind a draft citation.

#### Scenario: Citation lookup succeeds
- **WHEN** a reviewer requests inspection for a citation that exists in the corpus
- **THEN** the backend returns the matching source document and chunk details
- **AND** the response includes the cited excerpt and source location

#### Scenario: Citation lookup fails
- **WHEN** a reviewer requests inspection for a citation that does not exist in the corpus
- **THEN** the backend returns a not-found style error
- **AND** no empty inspection payload is returned

### Requirement: Workbench MUST expose citation inspection during review
The Speclist workbench MUST allow reviewers to inspect citations from retrieval results and draft sections.

#### Scenario: Reviewer inspects a draft citation
- **WHEN** a reviewer selects a citation while editing a draft
- **THEN** the workbench loads the cited source context
- **AND** shows the source title, location, section, and excerpt in a dedicated panel
