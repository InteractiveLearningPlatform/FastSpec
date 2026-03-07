## Why

OpenSpec-targeted export now uses typed templates, but FastSpec YAML export still emits one generic draft document shape. The next useful step is to make FastSpec export render more intentional YAML structures so the output is a better starting point for durable spec refinement.

## What Changes

- Add typed FastSpec YAML rendering for exported drafts.
- Render different YAML shapes depending on the draft content and export intent.
- Keep citations and summary metadata in the exported YAML and sidecar output.

## Capabilities

### New Capabilities
- `speclist-typed-fastspec-export`: render typed FastSpec YAML exports instead of a single generic draft document

### Modified Capabilities
- `draft-export`: FastSpec YAML export becomes typed and structured rather than generic draft serialization

## Impact

- Affected code: `apps/speclist-api`
- Affected behavior: FastSpec YAML export becomes a stronger starting point for durable spec authoring
- Affected tests: YAML export rendering coverage
