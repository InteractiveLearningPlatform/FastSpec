## Why

Speclist can already ingest sources, retrieve grounded context, and produce reviewable draft specs, but it still leaves accepted drafts stranded in the UI. The next useful slice is to export reviewed drafts into durable files so teams and agents can move from draft review into real FastSpec or OpenSpec artifacts without manual copy-paste.

## What Changes

- Add a draft export flow to Speclist.
- Export reviewed drafts into explicit target directories on disk.
- Support at least one OpenSpec markdown export format and one FastSpec YAML export format.
- Write machine-readable citation metadata alongside exported artifacts.
- Extend the React workbench with export controls after draft review.

## Capabilities

### New Capabilities
- `draft-export`: export reviewed Speclist drafts into durable OpenSpec and FastSpec files with citation metadata

### Modified Capabilities
- `speclist-workbench`: draft review now includes a concrete export action instead of stopping at on-screen inspection

## Impact

- Affected code: `apps/speclist-api`, `apps/speclist-web`
- Affected behavior: reviewed drafts can now become durable files in the local repo or another explicit output path
- Affected docs: Speclist product and local workflow documentation
- Affected tests: backend export tests and frontend build verification
