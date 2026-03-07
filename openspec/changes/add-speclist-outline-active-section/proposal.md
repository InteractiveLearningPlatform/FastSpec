## Why

Speclist reviewers can now navigate long drafts through the outline, but the outline still behaves like a plain list of links. Reviewers need a visible active-section marker so they can tell which part of the draft is currently in focus after jumping or editing.

## What Changes

- Track the active draft section during review.
- Highlight the active section in the draft outline.
- Keep active-section tracking aligned with outline navigation and local draft resets.

## Capabilities

### New Capabilities
- `speclist-outline-active-section`: show which draft section is currently active in the outline

### Modified Capabilities
- `speclist-draft-outline`: outline entries now communicate current focus state
- `speclist-workbench`: reviewers can orient themselves more quickly after navigating long drafts

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: outline navigation becomes stateful and easier to follow
- Affected tests: frontend build verification
