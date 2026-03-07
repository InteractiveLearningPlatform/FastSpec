## Why

Speclist reviewers can now reorder sections, but they still cannot quickly branch an existing section into a variation during review. Reviewers need a lightweight duplication action so they can copy a section and refine the copy without manual copy-paste.

## What Changes

- Add a section duplication action to the draft editor.
- Insert the duplicate next to the original section.
- Copy any existing review flag state to the duplicated section.

## Capabilities

### New Capabilities
- `speclist-section-duplicate`: duplicate a draft section during review

### Modified Capabilities
- `speclist-workbench`: reviewers can branch an existing section into a local variation

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can duplicate sections without recreating them manually
- Affected tests: frontend build verification
