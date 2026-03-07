## Context

The draft outline already derives from the current section list and uses explicit click navigation. That means active-section tracking can remain entirely in local UI state. The main requirement is to keep the marker predictable when users generate a new draft, reset review state, or jump through the outline.

## Goals / Non-Goals

**Goals:**
- Track a current active section in the draft review session.
- Reflect the active section in the outline.
- Update active state when a reviewer jumps from the outline or resets the draft.

**Non-Goals:**
- Track viewport position with scroll observers.
- Persist active-section state across reloads.
- Add keyboard navigation in this slice.

## Decisions

Represent the active section as a single section index in local component state.
Rationale: the outline is already index-aligned with the draft section list, so one index is enough.

Set the active section when outline navigation is used.
Rationale: explicit reviewer navigation is the clearest signal of intended focus.

Reset the active section to the first section when a new draft is generated or the draft is reset.
Rationale: those actions establish a fresh review baseline.

Clear the active state when no draft exists.
Rationale: the outline should not show stale focus after draft removal.

## Risks / Trade-offs

[The active marker may not follow manual scrolling] -> This slice intentionally follows explicit navigation and reset events only, which keeps the behavior simple and predictable.

[Index-based active state can drift if section operations are not handled carefully] -> The active index is updated only by draft lifecycle and outline navigation in this slice, avoiding broader coupling to every section edit action.
