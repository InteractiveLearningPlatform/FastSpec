## Why

The outline now supports filtering and active-section tracking, but reviewers still have to click entries one by one to move through a draft. A lightweight next/previous flow is the next useful refinement because it lets reviewers walk the current outline view without manually selecting every section.

## What Changes

- Add next and previous navigation controls for the current outline view.
- Make the controls respect active section state and active outline filters.
- Keep navigation local to the review surface without changing draft content.

## Capabilities

### New Capabilities
- `speclist-next-section-nav`: move forward or backward through the current outline view

### Modified Capabilities
- `speclist-draft-outline`: reviewers can step through outline entries sequentially
- `speclist-outline-filtering`: guided navigation follows the filtered outline view

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: review passes over long drafts become more efficient
- Affected tests: frontend build verification
