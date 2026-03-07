## Context

The current workbench shows editable drafts and diffs, but reviewers still have to keep unresolved issues in their head. Since draft review is intentionally in-memory, the next useful step is to attach lightweight review metadata directly to draft sections in the browser.

## Goals / Non-Goals

**Goals:**
- Add per-section review status and optional notes.
- Make flagged sections visible in the main draft editor and in one summary view.
- Keep the feature entirely client-side for now.

**Non-Goals:**
- Persist flags to the backend or exported artifacts.
- Add reviewer identity, workflow assignment, or approval gates.
- Build a full issue tracker inside the workbench.

## Decisions

Store review flags in local workbench state keyed by section index.
Rationale: the feature is for immediate review context and does not need backend persistence in this slice.

Use a small fixed status set: `ready`, `needs-work`, and `blocked`.
Rationale: these three states are enough to express the next action without adding process overhead.

Render both inline controls and a flagged-section summary panel.
Rationale: reviewers need local context while editing and a compact pre-export overview.

## Risks / Trade-offs

[Section-index keys may shift after reordering or removal] -> Keep the implementation simple now and accept index-based state because section reordering is not yet a primary workflow.

[Flags disappear on reload] -> Preserve the existing lightweight in-memory review model and defer persistence to a later slice if it proves necessary.
