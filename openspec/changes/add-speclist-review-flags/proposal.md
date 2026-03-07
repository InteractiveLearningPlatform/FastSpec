## Why

Speclist reviewers can now edit drafts and inspect diffs, citations, and source detail, but they still cannot mark sections that need follow-up before export. Lightweight review flags would let reviewers keep track of concerns without leaving the workbench.

## What Changes

- Add section-level review flags to the Speclist workbench.
- Support a small set of flag states and optional reviewer notes.
- Add a review summary panel that shows flagged sections before export.

## Capabilities

### New Capabilities
- `speclist-review-flags`: mark draft sections with lightweight review status before export

### Modified Capabilities
- `speclist-workbench`: draft review now includes explicit section-level follow-up markers

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can mark sections as ready, needs work, or blocked with notes
- Affected tests: frontend build verification
