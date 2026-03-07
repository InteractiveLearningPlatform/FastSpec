## Why

Speclist reviewers can add and remove draft sections, but they still cannot quickly rearrange section order once the review evolves. Reviewers need a lightweight reorder action so the draft structure can be refined without rebuilding sections manually.

## What Changes

- Add section move-up and move-down actions to the draft editor.
- Reorder section review state together with the sections.
- Keep draft diff and readiness views consistent after section moves.

## Capabilities

### New Capabilities
- `speclist-section-reorder`: rearrange draft sections during review

### Modified Capabilities
- `speclist-workbench`: draft structure can be refined by reordering existing sections

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can reorder sections without recreating them
- Affected tests: frontend build verification
