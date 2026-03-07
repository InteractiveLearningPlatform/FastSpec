## Context

The current workbench supports importing and indexing many sources, but the indexed source list is only a summary view. Since all source documents and their chunks are already stored in memory and exposed through the backend store, the simplest next step is to provide exact source lookup by id and render that context in a dedicated panel.

## Goals / Non-Goals

**Goals:**
- Resolve a source id to the full source document with metadata and chunks.
- Add a direct inspect action to indexed source cards.
- Show source metadata and chunk inventory in a review-friendly panel.

**Non-Goals:**
- Add source editing or deletion.
- Build pagination or virtualization for extremely large source documents.
- Introduce platform-scale search or storage changes from the separate ops track.

## Decisions

Use exact source-id lookup in the backend.
Rationale: the UI already has stable source ids from `/api/v1/sources`, so there is no need for an additional search layer.

Return the full `SourceDocument` shape for source detail inspection.
Rationale: the UI needs both metadata and the chunk list, and the existing domain object already captures that.

Keep the detail panel read-only and separate from citation inspection.
Rationale: source detail view is corpus exploration, while citation inspection is evidence verification from retrieval and draft review.

## Risks / Trade-offs

[Large source documents may produce long panels] -> Keep the view simple now and defer chunk collapsing or pagination to a later slice if it becomes necessary.

[Exact id lookup assumes source cards are current] -> Use ids from the current `/sources` payload and return a clear error if the source is missing.
