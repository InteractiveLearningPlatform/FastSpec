## 1. Repository and service scaffolding

- [x] 1.1 Add `speclist-api` and `speclist-web` application roots with docs for local development
- [x] 1.2 Define the backend hexagonal module structure, service boundaries, and configuration model
- [x] 1.3 Document how imported docs, existing specs, and generated drafts fit into the repo workflow

## 2. Ingestion and corpus pipeline

- [x] 2.1 Implement DOCX ingestion with normalized chunk extraction and source metadata
- [x] 2.2 Implement Confluence ingestion with page import and provenance tracking
- [x] 2.3 Persist normalized document chunks and structure metadata in a corpus store

## 3. Retrieval and draft generation

- [x] 3.1 Index existing FastSpec/OpenSpec specs together with imported document chunks
- [x] 3.2 Implement spec-oriented retrieval bundles with citations and compact context output
- [x] 3.3 Implement draft generation endpoints that create reviewable spec candidates from retrieved context

## 4. Frontend workbench

- [x] 4.1 Build React flows for document import, source browsing, and retrieval search
- [x] 4.2 Add draft review UI that shows generated content with supporting evidence
- [x] 4.3 Add end-to-end examples or fixtures that demonstrate DOCX/Confluence to spec drafting
