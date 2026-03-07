## Why

The draft outline now supports navigation and active-state orientation, but it still becomes noisy on larger drafts. Reviewers need a lightweight way to narrow the outline to the sections they are looking for or the sections that still need work.

## What Changes

- Add a text filter for outline entries.
- Add a quick filter that shows only non-ready sections.
- Keep filtering local to the outline without mutating the underlying draft.

## Capabilities

### New Capabilities
- `speclist-outline-filtering`: filter outline entries by text or review state

### Modified Capabilities
- `speclist-draft-outline`: reviewers can narrow long outlines to a smaller working set
- `speclist-review-flags`: review-state metadata can drive outline filtering

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: long draft navigation becomes easier when the outline is dense
- Affected tests: frontend build verification
