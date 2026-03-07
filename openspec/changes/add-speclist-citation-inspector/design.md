## Context

The current workbench shows citations in both retrieval results and draft sections, but it does not provide a way to reopen the actual chunk behind a citation. Since citations are persisted as stable strings on chunks, the simplest useful flow is a backend lookup that resolves one citation string back to its source document and chunk.

## Goals / Non-Goals

**Goals:**
- Resolve a citation string to the matching source document and chunk.
- Make citations clickable from retrieval results and draft review.
- Present the source title, location, section, excerpt, and metadata in one place.

**Non-Goals:**
- Build full corpus browsing or multi-result citation faceting.
- Add persistent reviewer annotations on citations.
- Introduce vector search or marketplace ranking work from the platform-ops change.

## Decisions

Use a dedicated inspection endpoint keyed by citation string.
Rationale: draft citations are currently stored as strings, so citation lookup should work even after reviewers edit section bodies or when the relevant search result is no longer on screen.

Return one exact citation match with source and chunk context.
Rationale: citations are generated as specific chunk labels, so exact-match lookup is the right default for this slice.

Expose the same inspect action from retrieval results and draft citations.
Rationale: reviewers should use one mental model for evidence inspection across the workbench.

Keep the frontend panel read-only.
Rationale: the goal is evidence verification, not source editing.

## Risks / Trade-offs

[Manually edited citations may not resolve] -> Return a clear backend error so reviewers can fix the citation text.

[Exact string matching may be strict] -> Prefer stable deterministic lookup now and defer fuzzy citation matching to later work if needed.
