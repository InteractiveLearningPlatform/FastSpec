## ADDED Requirements

### Requirement: Platform MUST provide ML, RAG, and NLP pipelines for specs
The platform MUST provide ML, RAG, and NLP pipelines for ingesting, normalizing, enriching, ranking, and generating specs.

#### Scenario: Ingested source enriched
- **WHEN** a new source is ingested into Speclist
- **THEN** the platform runs normalization and enrichment steps
- **AND** the resulting indexed asset is suitable for retrieval and generation workflows

### Requirement: Pipelines MUST support code and IR-aware processing
The platform MUST support indexing pipelines for source code, GitHub markdown, and language-oriented intermediate representations in addition to prose documents.

#### Scenario: Code and IR indexed
- **WHEN** the platform ingests a repository or code-oriented source
- **THEN** it can index code, markdown, and supported intermediate representations
- **AND** those indexed assets are available to retrieval and generation workflows
