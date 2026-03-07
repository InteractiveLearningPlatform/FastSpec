## Context

The current workbench keeps both the original generated draft and the editable draft, but there is no action that uses that stored original snapshot to reset the review. Since the original draft already exists in client state, the simplest useful change is a local reset action.

## Goals / Non-Goals

**Goals:**
- Reset the editable draft to the original generated snapshot.
- Clear review-only state that no longer applies after reset.
- Keep the feature entirely client-side.

**Non-Goals:**
- Add undo/redo history.
- Ask the backend for a new draft during reset.
- Persist reset history.

## Decisions

Restore the editable draft from the preserved original snapshot.
Rationale: this is already the canonical generated baseline used by the diff panel.

Clear review flags and export result state during reset.
Rationale: those states are attached to the edited review session and should not survive a reset.

Keep the original generated snapshot unchanged after reset.
Rationale: the reset target should remain stable until the reviewer generates a new draft.

## Risks / Trade-offs

[Reset is destructive for in-memory edits] -> Keep the action explicit in the review surface and scoped to the current draft only.

[Some teams may want partial reset] -> Start with full reset to baseline and defer finer controls if needed later.
