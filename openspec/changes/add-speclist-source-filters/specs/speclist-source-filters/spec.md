## ADDED Requirements

### Requirement: Retrieval MUST support source filters
The system MUST allow retrieval requests to limit candidate sources by source kind, source origin, and location substring before ranking chunks.

#### Scenario: Retrieval filtered to repository specs
- **WHEN** an operator searches with the source origin filter set to `repository`
- **THEN** only repository-indexed spec sources contribute chunks to the retrieval bundle
- **AND** imported DOCX or Confluence sources are excluded

#### Scenario: Retrieval filtered by location text
- **WHEN** an operator searches with a location substring filter
- **THEN** only sources whose location contains that substring contribute chunks
- **AND** ranking is applied only to the filtered source set

### Requirement: Draft generation MUST reuse retrieval filters
The system MUST apply the same source filters during draft generation so the resulting draft is grounded in the intended subset of sources.

#### Scenario: Draft generated from imported documents only
- **WHEN** an operator drafts a spec with source origin set to `imported`
- **THEN** the generated draft uses retrieval results from imported sources only
- **AND** the reported source count reflects the filtered result set

### Requirement: Workbench MUST expose retrieval filters
The Speclist workbench MUST provide simple controls for source kind, source origin, and location filtering when searching and drafting.

#### Scenario: Operator narrows search from the UI
- **WHEN** an operator selects one or more source kinds and enters a location filter in the workbench
- **THEN** the frontend sends those filters with search and draft requests
- **AND** the applied filters stay visible while reviewing retrieval results
