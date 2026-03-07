## Context

The current workbench supports editing, adding, removing, diffing, flagging, readiness checks, and reset, but section order is still fixed unless the reviewer manually recreates sections. Since both the section list and review flag state already live in client memory, the cheapest useful improvement is an index-based reorder action.

## Goals / Non-Goals

**Goals:**
- Add simple move-up and move-down actions for draft sections.
- Keep section review flags aligned with their moved sections.
- Preserve compatibility with the existing diff and readiness logic.

**Non-Goals:**
- Add drag-and-drop or arbitrary position input.
- Reorder source citations independently from sections.
- Introduce server-side ordering state.

## Decisions

Use explicit up/down buttons per section.
Rationale: this is straightforward, low-risk, and matches the current lightweight workbench approach.

Reorder review flags together with the section list.
Rationale: flags are review metadata for a specific section and must move with that section.

Keep diff comparison index-based even after reorder.
Rationale: the existing diff model is intentionally simple; reordering support should not expand into complex matching in this slice.

## Risks / Trade-offs

[Index-based diff output may look noisier after reordering] -> Accept this for now and defer richer matching if reordering becomes the dominant workflow.

[Repeated button clicks are less fluid than drag-and-drop] -> Prefer the simpler control surface first.
