## ADDED Requirements

### Requirement: System MUST retrieve spec-authoring context from both imported documentation and existing specs
The system MUST retrieve grounded context bundles that combine imported source material with relevant existing FastSpec or OpenSpec specs.

#### Scenario: Retrieval returns mixed context bundle
- **WHEN** a user or agent searches for context to create a new spec
- **THEN** the system returns relevant existing specs and imported documentation excerpts
- **AND** the response is shaped as a compact bundle for spec-authoring workflows

### Requirement: Retrieval MUST support draft generation for new specs
The system MUST provide a workflow that uses retrieved context to draft new specs while preserving traceable evidence for generated content.

#### Scenario: Draft spec generated from retrieval context
- **WHEN** a user requests a new draft spec from retrieved context
- **THEN** the system produces a draft artifact candidate
- **AND** the draft includes references to the supporting retrieved sources

### Requirement: Retrieval MUST remain simple and efficient
The retrieval workflow MUST remain fast enough for interactive use and MUST avoid returning unnecessarily large raw documents when a compact context bundle is available.

#### Scenario: Compact result returned
- **WHEN** the query matches a large imported document
- **THEN** the system returns the most relevant structured excerpts instead of the full document
- **AND** the response remains suitable for interactive frontend use
