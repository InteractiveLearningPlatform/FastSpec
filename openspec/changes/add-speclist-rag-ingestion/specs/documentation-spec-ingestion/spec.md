## ADDED Requirements

### Requirement: System MUST ingest DOCX and Confluence documentation into a spec-ready corpus
The system MUST ingest documentation from uploaded DOCX files and configured Confluence pages, normalize the content, and persist it as a spec-ready corpus with source references.

#### Scenario: DOCX document ingested
- **WHEN** an operator uploads a DOCX file to `speclist`
- **THEN** the backend stores the source document metadata and extracted text
- **AND** it creates normalized chunks suitable for retrieval and spec drafting

#### Scenario: Confluence page ingested
- **WHEN** an operator imports a Confluence page or page tree
- **THEN** the backend fetches the selected page content and metadata
- **AND** it stores normalized retrieval chunks with a traceable Confluence source reference

### Requirement: Ingestion MUST preserve provenance and citations
The system MUST preserve enough provenance metadata for every imported chunk so generated specs can cite the original source material.

#### Scenario: Retrieved chunk includes source evidence
- **WHEN** a retrieval result contains imported documentation
- **THEN** each returned chunk includes a stable source identifier
- **AND** the result includes enough metadata to link back to the original DOCX or Confluence source

### Requirement: Ingestion MUST extract spec-oriented structure
The system MUST extract headings, sections, and other useful structure from imported documents so retrieval can prefer material that maps cleanly to spec authoring.

#### Scenario: Structured sections available after import
- **WHEN** a document is successfully imported
- **THEN** normalized chunks retain section-level context
- **AND** the corpus stores structure metadata that can be used during retrieval and draft generation
