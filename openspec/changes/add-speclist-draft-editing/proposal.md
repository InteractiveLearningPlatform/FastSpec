## Why

Speclist can retrieve grounded context and generate draft specs, but reviewers still have to accept the generated text as-is before export. The next practical step is to let operators refine draft titles, summaries, and section content directly in the workbench so exported artifacts reflect human review rather than raw generation output.

## What Changes

- Add draft editing controls to the Speclist workbench.
- Allow reviewers to edit the draft title, summary, section headings, section bodies, and citations before export.
- Validate and normalize edited drafts on export so malformed reviewer changes do not write broken artifacts.

## Capabilities

### New Capabilities
- `speclist-draft-editing`: review and refine generated drafts before export

### Modified Capabilities
- `speclist-workbench`: draft review becomes editable instead of read-only
- `draft-export`: export validates edited draft structure before writing files

## Impact

- Affected code: `apps/speclist-web`, `apps/speclist-api`
- Affected behavior: reviewers can adjust generated draft content before exporting it
- Affected tests: backend export validation coverage
