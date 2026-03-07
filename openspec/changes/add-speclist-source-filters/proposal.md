## Why

Speclist can already ingest documents, index repository specs, and retrieve grounded chunks, but retrieval still searches the whole corpus every time. Operators need a simple way to narrow searches to imported docs, repository specs, or specific source paths so draft generation stays focused without crossing into the separate platform-ops work.

## What Changes

- Add source filters to Speclist retrieval and draft generation.
- Support filtering by source kind, source origin, and location substring.
- Add workbench controls for applying those filters before search and draft generation.

## Capabilities

### New Capabilities
- `speclist-source-filters`: narrow retrieval and draft generation to selected source slices

### Modified Capabilities
- `spec-rag-retrieval`: retrieval requests can apply source filters before ranking chunks
- `speclist-workbench`: operators can configure retrieval filters from the UI

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected behavior: retrieval and draft generation become more targeted across mixed source corpora
- Affected tests: backend service and HTTP coverage for filtered search behavior
