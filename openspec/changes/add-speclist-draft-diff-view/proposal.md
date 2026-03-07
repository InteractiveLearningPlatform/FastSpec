## Why

Speclist now supports editable drafts and preset-aware generation, but reviewers still cannot quickly see what changed after they start editing. A diff view inside the workbench would make review safer by showing how the current draft diverges from the originally generated draft before export.

## What Changes

- Capture the original generated draft in the workbench.
- Render a draft diff view that compares the original and current draft.
- Highlight changes across title, summary, and section content before export.

## Capabilities

### New Capabilities
- `speclist-draft-diff-view`: compare the current edited draft against the original generated draft

### Modified Capabilities
- `speclist-workbench`: draft review now includes explicit change inspection before export

## Impact

- Affected code: `apps/speclist-web`
- Affected behavior: reviewers can inspect edits before exporting a modified draft
- Affected tests: frontend build verification
