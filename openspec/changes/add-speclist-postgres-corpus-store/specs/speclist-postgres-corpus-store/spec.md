## ADDED Requirements

### Requirement: Speclist MUST support PostgreSQL-backed corpus persistence
Speclist MUST support storing imported and indexed source documents in PostgreSQL through the existing corpus store contract.

#### Scenario: Source document saved to PostgreSQL
- **WHEN** the service runs with PostgreSQL corpus storage enabled and saves a source document
- **THEN** the document and its chunks are persisted in PostgreSQL
- **AND** later reads return the same document through the `CorpusStore` interface

### Requirement: Speclist MUST keep file storage as an explicit fallback
Speclist MUST allow operators to keep using file-backed corpus storage when PostgreSQL is not selected.

#### Scenario: File store remains selected
- **WHEN** the service runs without PostgreSQL corpus storage enabled
- **THEN** it continues to use file-backed corpus persistence
- **AND** existing file-based behavior remains available for local development

### Requirement: PostgreSQL corpus storage MUST bootstrap its schema
The PostgreSQL corpus store MUST create the tables it requires before serving reads or writes.

#### Scenario: Empty database initialized
- **WHEN** the PostgreSQL corpus store starts against an empty database
- **THEN** it creates the required tables for source documents and chunks
- **AND** the store becomes ready without a separate manual migration step
