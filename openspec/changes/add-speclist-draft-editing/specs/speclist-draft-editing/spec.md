## ADDED Requirements

### Requirement: Workbench MUST allow editing generated drafts
The Speclist workbench MUST allow reviewers to modify generated draft content before export.

#### Scenario: Reviewer edits generated draft content
- **WHEN** a draft has been generated in the workbench
- **THEN** the reviewer can edit the draft title, summary, section headings, section bodies, and citations
- **AND** those edits are retained in the in-memory draft state until export or replacement

### Requirement: Workbench MUST support simple section reshaping
The Speclist workbench MUST allow reviewers to add and remove draft sections before export.

#### Scenario: Reviewer adds a missing section
- **WHEN** the generated draft needs another section
- **THEN** the reviewer can append a new section in the workbench
- **AND** the new section is included in the exported draft payload

### Requirement: Export MUST validate edited drafts
The backend MUST validate reviewer-edited drafts before writing exported artifacts.

#### Scenario: Export rejects invalid edited draft
- **WHEN** an edited draft contains a blank title or an empty section heading/body
- **THEN** export fails with a validation error
- **AND** no primary artifact or sidecar file is written
