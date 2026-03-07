## Why

FastSpec currently starts from hand-authored YAML, which is too limiting when teams already have system knowledge trapped in DOCX files and Confluence pages. The next useful product slice is a lightweight service that can ingest those documents, build retrieval-quality spec context, and help agents draft new specs from grounded source material instead of inventing structure.

## What Changes

- Add a `speclist` service with a Go backend and React frontend.
- Ingest documentation from DOCX files and Confluence pages into a normalized spec-oriented corpus.
- Store document chunks, extracted metadata, and source citations so generated specs remain grounded.
- Provide spec-focused retrieval and RAG workflows that return relevant existing specs plus source evidence for drafting new specs.
- Expose a small UI for uploading sources, browsing indexed specs, searching context, and creating draft specs.
- Define the backend as a microservice using hexagonal architecture to keep adapters separate from core ingestion and retrieval logic.

## Capabilities

### New Capabilities
- `documentation-spec-ingestion`: ingest DOCX and Confluence content into a normalized spec-ready corpus with traceable citations
- `spec-rag-retrieval`: retrieve existing specs and imported documentation context for grounded spec drafting
- `speclist-workbench`: provide the `speclist` backend and frontend surfaces for ingestion, search, and draft generation workflows

### Modified Capabilities
- `repository-foundation`: the repository layout and examples must support a Go backend plus React frontend service alongside the Rust FastSpec tooling

## Impact

- Affected code: new `apps/speclist-api` and `apps/speclist-web` applications, plus shared docs and local dev orchestration
- Affected architecture: introduce a Go hexagonal backend for ingestion/retrieval and a React UI for operator workflows
- Affected dependencies: DOCX parsing, Confluence API integration, search/vector indexing, and frontend build tooling
- Affected docs: product architecture, local development, and example ingestion/generation workflows
