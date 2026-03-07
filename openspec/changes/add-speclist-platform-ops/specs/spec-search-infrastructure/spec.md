## ADDED Requirements

### Requirement: Search architecture MUST remain DB-agnostic with a concrete default stack
The platform MUST keep search and indexing interfaces DB-agnostic while supporting PostgreSQL, ClickHouse, and Valkey as the default operational stack.

#### Scenario: Storage backends connected through stable ports
- **WHEN** contributors change one storage implementation
- **THEN** the application domain contracts remain stable
- **AND** the default deployment can still use PostgreSQL, ClickHouse, and Valkey together

### Requirement: Search stack MUST support scalable vector and hybrid retrieval
The platform MUST support scalable vector and hybrid retrieval for specs, markdown, code, and other indexed technical assets.

#### Scenario: Hybrid search returns relevant technical context
- **WHEN** a query spans structured specs and technical source material
- **THEN** the search layer returns relevant hybrid results
- **AND** the platform can rank them for spec-authoring workflows

### Requirement: Initial vector retrieval baseline MUST support code-oriented search
The platform MUST choose an initial vector retrieval baseline that supports code-oriented semantic search and scalable retrieval.

#### Scenario: Code search baseline documented
- **WHEN** contributors review the platform design
- **THEN** the documented baseline includes a code-oriented vector retrieval engine
- **AND** the decision explains how it supports code, markdown, and hybrid retrieval
