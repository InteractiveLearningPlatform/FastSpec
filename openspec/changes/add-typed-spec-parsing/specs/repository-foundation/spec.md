## MODIFIED Requirements

### Requirement: Minimal Executable Rust Scaffold Exists
The repository MUST include a minimal Rust workspace that can parse and inspect FastSpec example documents.

#### Scenario: Contributor validates example specs
- **WHEN** a contributor runs the workspace tests or CLI
- **THEN** the workspace can detect supported FastSpec document kinds
- **AND** it can parse example documents into typed runtime structures
- **AND** it can summarize or inspect the example tree without external services
