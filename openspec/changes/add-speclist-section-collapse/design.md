## Context

The current workbench keeps every section fully expanded. That works for short drafts, but longer review sessions now include section editing, flags, duplication, and reordering, which makes the draft panel visually noisy. Collapse state can stay entirely client-side alongside the rest of the review-only UI state.

## Goals / Non-Goals

**Goals:**
- Let reviewers collapse and expand individual draft sections.
- Keep the section heading and controls visible while collapsed.
- Preserve collapse state across local section review actions.

**Non-Goals:**
- Persist collapse state on the backend.
- Add nested outlines or multi-level folding.
- Introduce bulk collapse controls in this slice.

## Decisions

Store collapsed section state as an index-aligned local list in the workbench.
Rationale: section review flags already follow this pattern, and the feature stays local to the current draft session.

Render collapsed sections as compact cards that keep heading and controls visible.
Rationale: reviewers need to preserve orientation in the draft while reducing visual noise.

Preserve collapse state when sections move or duplicate by transforming the collapse-state list alongside the section list.
Rationale: collapse is part of the review session state and should track the section structure predictably.

Reset collapse state when the draft is reset to the generated baseline.
Rationale: reset should restore a clean starting point for review.

## Risks / Trade-offs

[Tracking collapse state by position can drift if list transformations are inconsistent] -> Reuse the same index-based helper pattern already used for section review flags.

[Collapsed sections may hide issues from reviewers] -> The heading, flags, and controls stay visible so collapsed sections remain discoverable.
