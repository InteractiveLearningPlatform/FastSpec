## Why

Speclist now supports citation inspection, but the indexed source list is still shallow and only shows high-level labels. Reviewers need to open a source directly to inspect its metadata and chunk inventory without first finding a citation that points into it.

## What Changes

- Add backend source detail lookup by source id.
- Expose source detail lookup through the Speclist API.
- Add source detail actions and a dedicated detail panel in the workbench.

## Capabilities

### New Capabilities
- `speclist-source-detail-view`: inspect a source document and its chunks from the indexed source list

### Modified Capabilities
- `speclist-workbench`: indexed sources become navigable review entries instead of static labels

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected behavior: operators can inspect source-level metadata and chunk context directly from the corpus list
- Affected tests: backend service and HTTP coverage for source detail lookup
