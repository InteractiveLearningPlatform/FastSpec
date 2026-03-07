## Context

The outline is already a derived rendering of the current draft sections. Filtering can stay entirely in that derived layer, which means no section content or export behavior needs to change. The most useful baseline is a simple text match and a review-state shortcut for unresolved sections.

## Goals / Non-Goals

**Goals:**
- Filter outline entries by text.
- Filter outline entries to only non-ready sections.
- Keep navigation working against the original section index.

**Non-Goals:**
- Filter the draft editor itself.
- Add advanced query syntax.
- Persist filters across reloads.

## Decisions

Store outline filters as local UI state near the draft outline.
Rationale: the filters only affect the outline presentation, not the draft content.

Support case-insensitive text matching against section headings.
Rationale: headings are the main navigational signal in the outline.

Support a boolean `only non-ready` filter based on review status and notes.
Rationale: reviewers often need to jump between unresolved sections rather than all sections.

Preserve navigation by carrying the original section index through filtered outline entries.
Rationale: the filtered outline is a view over the real draft, not a new section ordering.

## Risks / Trade-offs

[Filtered outlines can hide the active section] -> Filtering is explicit and reversible, and the underlying draft remains unchanged.

[Review-state filtering based only on status could miss noted sections] -> Treat non-empty review notes as non-ready for filtering purposes.
