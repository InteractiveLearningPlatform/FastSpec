## Context

The current workbench keeps only the live editable draft state. Since draft editing is intentionally in-memory, the cheapest way to support review diffs is to also retain a snapshot of the originally generated draft and compare it client-side.

## Goals / Non-Goals

**Goals:**
- Preserve the original generated draft separately from the editable copy.
- Show differences in title, summary, preset, and sections.
- Keep the comparison lightweight and visible before export.

**Non-Goals:**
- Add server-side diffing or draft version persistence.
- Build a line-by-line text diff engine.
- Add collaborative review or comment threads.

## Decisions

Capture the original generated draft snapshot when `/drafts` returns.
Rationale: that keeps the implementation local to the workbench and avoids expanding the backend contract.

Show a structured comparison rather than a textual patch.
Rationale: reviewers care about which draft fields and sections changed, not about a raw unified diff.

Treat section comparison by index with explicit added/removed markers.
Rationale: section reordering is not currently a major workflow, and index-based comparison is enough for this slice.

## Risks / Trade-offs

[Index-based section comparison is approximate after heavy reshaping] -> Keep the view simple now and defer more advanced matching if the workflow demands it later.

[Large drafts may produce long diff panels] -> Reuse the existing workbench panel layout and keep each diff block compact.
