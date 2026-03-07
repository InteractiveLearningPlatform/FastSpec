## Context

The review surface already keeps all section data in one client-side draft object. That makes an outline purely a derived view of the current section list. The only extra behavior needed is navigation to a section, plus reopening the section if it is currently collapsed.

## Goals / Non-Goals

**Goals:**
- Show a compact ordered outline for the current draft.
- Let reviewers jump from the outline to a section.
- Expand the targeted section if it is collapsed.

**Non-Goals:**
- Persist outline state.
- Add nested hierarchy inference.
- Add drag-and-drop reordering from the outline.

## Decisions

Render the outline inside the draft review panel rather than as a separate app region.
Rationale: the feature is navigation for the active draft, not a new global surface.

Use section refs to scroll the draft section into view.
Rationale: this keeps navigation lightweight without introducing routing or URL fragments.

Clear the collapsed state for the target section before scrolling.
Rationale: navigation should land on editable content, not a still-collapsed card.

Show review status and collapsed state in the outline rows.
Rationale: reviewers should be able to scan section state before jumping.

## Risks / Trade-offs

[Outline items can become stale if they do not follow the current draft order] -> Build the outline directly from the current in-memory section list on every render.

[Automatic expansion may surprise reviewers who intended to keep sections folded] -> Only the targeted section is reopened, and the rest of the collapse state remains untouched.
