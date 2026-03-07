## Why

Speclist drafts are getting denser as reviewers edit, reorder, and duplicate sections. Reviewers need a lightweight way to collapse sections they are not actively editing so they can focus on one part of the draft without losing the overall structure.

## What Changes

- Add a collapse and expand action for each draft section in the Speclist workbench.
- Keep collapsed sections visible by heading so reviewers retain structural context.
- Preserve collapse state during in-memory review actions such as reordering, duplication, and reset.

## Capabilities

### New Capabilities
- `speclist-section-collapse`: collapse and expand draft sections during review

### Modified Capabilities
- `speclist-workbench`: reviewers can compress inactive sections while editing a draft

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: long drafts become easier to review without removing sections from view
- Affected tests: frontend build verification
