## ADDED Requirements

### Requirement: Speclist MUST provide a backend and frontend workbench
The `speclist` product MUST provide a Go backend service and a React frontend so contributors can manage ingestion, retrieval, and draft-generation workflows.

#### Scenario: Operator opens workbench
- **WHEN** an operator opens the `speclist` frontend
- **THEN** the UI can reach the backend service
- **AND** the operator can access ingestion, search, and draft-generation actions from one workbench

### Requirement: Backend MUST follow hexagonal architecture
The `speclist` backend MUST separate core ingestion, retrieval, and generation logic from infrastructure adapters.

#### Scenario: Source adapter replaced without domain rewrite
- **WHEN** a contributor swaps one ingestion or storage adapter for another
- **THEN** the core application logic remains unchanged
- **AND** the integration occurs through explicit ports and adapters

### Requirement: Workbench MUST expose reviewable draft output
The frontend MUST let an operator inspect retrieved context and review generated draft specs before accepting them as durable artifacts.

#### Scenario: Draft reviewed in UI
- **WHEN** the backend returns a draft spec candidate
- **THEN** the frontend shows the draft content together with its supporting context
- **AND** the operator can review the draft before exporting or applying it
