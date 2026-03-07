## ADDED Requirements

### Requirement: Retrieval MUST scale beyond document-only search
The retrieval workflow MUST support scalable search across structured specs, markdown, code, and IR-oriented indexed assets.

#### Scenario: Technical search spans multiple indexed asset classes
- **WHEN** a query requires spec, markdown, and code context together
- **THEN** the retrieval layer can search across those indexed asset classes
- **AND** the response remains shaped for spec-authoring or marketplace workflows
