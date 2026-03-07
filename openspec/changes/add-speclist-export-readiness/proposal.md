## Why

Speclist reviewers can now inspect diffs and set section-level review flags, but export still behaves like a blind final step. Reviewers need a compact readiness check that tells them whether the draft is actually ready to export or which blockers still remain.

## What Changes

- Add a client-side export readiness check to the Speclist workbench.
- Summarize blocking and warning conditions before export.
- Surface readiness state directly next to the export controls.

## Capabilities

### New Capabilities
- `speclist-export-readiness`: evaluate whether the current draft is ready for export

### Modified Capabilities
- `speclist-workbench`: export review now includes an explicit readiness summary

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can see blockers and warnings before exporting
- Affected tests: frontend build verification
