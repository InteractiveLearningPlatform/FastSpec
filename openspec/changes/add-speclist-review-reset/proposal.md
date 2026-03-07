## Why

Speclist reviewers can now make substantial in-workbench changes to drafts, but there is still no quick way to discard that review state and return to the original generated draft. Reviewers need a simple reset action when they want to restart the review from the generated baseline.

## What Changes

- Add a one-step review reset action to the workbench.
- Restore the current draft back to the original generated snapshot.
- Clear review flags and export results when a reset is performed.

## Capabilities

### New Capabilities
- `speclist-review-reset`: reset the current review state back to the generated draft

### Modified Capabilities
- `speclist-workbench`: reviewers can restart draft review without generating a new draft

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can discard local review edits and return to the generated baseline
- Affected tests: frontend build verification
