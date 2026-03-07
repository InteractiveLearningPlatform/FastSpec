## Why

Section-level collapse made long drafts easier to review, but reviewers still have to toggle sections one by one when they want to focus on only a small part of the draft. Bulk controls are the next efficient step because they make it practical to compress or reopen the entire draft in one action.

## What Changes

- Add bulk `Collapse all` and `Expand all` controls to the draft review surface.
- Make bulk controls operate only on the current in-memory draft sections.
- Keep the controls aligned with existing local review state and draft reset behavior.

## Capabilities

### New Capabilities
- `speclist-collapse-all`: collapse or expand every draft section in one action

### Modified Capabilities
- `speclist-section-collapse`: reviewers can manage collapse state for the whole draft, not only one section at a time
- `speclist-workbench`: the draft review surface supports global structure focus controls

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can collapse or expand the whole draft without repeated per-section toggles
- Affected tests: frontend build verification
