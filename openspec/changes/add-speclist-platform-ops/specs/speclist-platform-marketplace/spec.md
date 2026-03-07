## ADDED Requirements

### Requirement: Speclist MUST operate as a spec marketplace platform
Speclist MUST support platform behaviors for publishing, discovering, ranking, and reusing specs rather than acting only as a local drafting tool.

#### Scenario: Operator searches marketplace
- **WHEN** an operator searches for specs in Speclist
- **THEN** the platform returns reusable specs and related assets from the indexed catalog
- **AND** the results support discovery beyond a single local project

### Requirement: Platform MUST support multiple spec source types
The marketplace platform MUST support structured specs, markdown documents, and other indexed technical sources as searchable marketplace assets.

#### Scenario: Mixed asset types shown
- **WHEN** a marketplace query matches multiple asset types
- **THEN** the platform can return structured specs together with related documentation or technical context
- **AND** the response identifies the source type for each result
