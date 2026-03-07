## Context

Speclist currently searches all stored sources and ranks chunks with a simple term-matching scorer. That works for small corpora, but once imported docs and repository specs live together, users need to constrain the search space to avoid drafting from the wrong source class.

## Goals / Non-Goals

**Goals:**
- Add a compact filter model for retrieval and draft generation.
- Support kind-based, origin-based, and location-substring filtering.
- Keep the backend domain and HTTP contracts explicit and easy to extend.
- Expose the filters in the React workbench without expanding into advanced search infrastructure.

**Non-Goals:**
- Implement platform-scale faceting, vector filtering, or marketplace ranking.
- Add database-backed indexing or ops/security infrastructure from the platform-ops change.
- Redesign the scoring model beyond filtering the candidate source set.

## Decisions

Introduce an explicit `RetrievalFilter` value object in the backend domain and reuse it for both search and draft generation.
Rationale: the same candidate-selection rules should apply whether the user wants retrieval results or a generated draft.

Model origin as a stable enum-like string with `imported` and `repository`.
Rationale: users care about whether sources came from uploaded docs or indexed specs, and that distinction can be derived from existing source kinds without new persistence requirements.

Apply filters at the source-document level before chunk scoring.
Rationale: this keeps the implementation simple and avoids wasting work scoring chunks from excluded sources.

Keep the UI to a small set of controls: source kind checkboxes, an origin selector, and a location substring input.
Rationale: that gives operators practical narrowing without pulling the workbench toward a full search-builder interface.

## Risks / Trade-offs

[Too many filters could hide relevant sources] -> Keep defaults empty so unfiltered behavior remains unchanged.

[Origin inference could be too coarse] -> Use the current split of spec vs imported sources and defer finer provenance categories to later work.

[Draft output could surprise users if filters differ from the last search] -> Reuse the same filter state for both search and draft actions in the UI.
