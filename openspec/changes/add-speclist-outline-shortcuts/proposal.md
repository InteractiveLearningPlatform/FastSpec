## Why

The outline now supports previous and next buttons, but long review passes still require repeated pointer movement back to the outline controls. Keyboard shortcuts are the next practical step because they speed up review without changing the underlying navigation model.

## What Changes

- Add keyboard shortcuts for previous and next outline navigation.
- Keep shortcuts scoped to active draft review sessions.
- Disable shortcuts while the reviewer is typing in form fields.

## Capabilities

### New Capabilities
- `speclist-outline-shortcuts`: navigate the draft outline with keyboard shortcuts

### Modified Capabilities
- `speclist-next-section-nav`: sequential navigation can be triggered from the keyboard
- `speclist-workbench`: reviewers can perform faster navigation during long review passes

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: review navigation becomes faster without requiring mouse interaction for every step
- Affected tests: frontend build verification
