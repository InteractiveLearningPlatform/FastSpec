## Why

Speclist draft generation still starts every review from the same generic section structure. Reviewers need a few intentional starting points so proposal-shaped work, design work, and requirement-heavy work do not all begin from the same scaffold.

## What Changes

- Add typed draft presets to the Speclist backend and draft payload.
- Support `general`, `proposal`, `design`, and `requirements` presets during draft generation.
- Add a preset selector to the workbench and preserve the selected preset on the generated draft.

## Capabilities

### New Capabilities
- `speclist-draft-presets`: generate drafts from preset review structures

### Modified Capabilities
- `speclist-workbench`: draft generation can start from different review-oriented structures
- `draft-export`: exported draft metadata preserves the selected preset

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected behavior: proposal, design, and requirements work can start from more intentional grounded section sets
- Affected tests: backend service and HTTP coverage for preset-aware draft generation
